import {ByteBuffer} from 'flatbuffers';
import {Connection} from '../../aos/network/www/proxy';
import {LocalizerOutput} from '../../frc971/control_loops/drivetrain/localization/localizer_output_generated';
import {RejectionReason} from '../localizer/status_generated';
import {Status as DrivetrainStatus} from '../../frc971/control_loops/drivetrain/drivetrain_status_generated';
import {Visualization, TargetEstimateDebug} from '../localizer/visualization_generated';

import {FIELD_LENGTH, FIELD_WIDTH, FT_TO_M, IN_TO_M} from './constants';

// (0,0) is field center, +X is toward red DS
const FIELD_SIDE_Y = FIELD_WIDTH / 2;
const FIELD_EDGE_X = FIELD_LENGTH / 2;

const ROBOT_WIDTH = 25 * IN_TO_M;
const ROBOT_LENGTH = 32 * IN_TO_M;

const PI_COLORS = ['#ff00ff', '#ffff00', '#00ffff', '#ffa500'];
const PIS = ['pi1', 'pi2', 'pi3', 'pi4'];

export class FieldHandler {
  private canvas = document.createElement('canvas');
  private localizerOutput: LocalizerOutput|null = null;
  private drivetrainStatus: DrivetrainStatus|null = null;

  // Image information indexed by timestamp (seconds since the epoch), so that
  // we can stop displaying images after a certain amount of time.
  private localizerImageMatches = new Map<number, Visualization>();
  private x: HTMLElement = (document.getElementById('x') as HTMLElement);
  private y: HTMLElement = (document.getElementById('y') as HTMLElement);
  private theta: HTMLElement =
      (document.getElementById('theta') as HTMLElement);
  private imagesAcceptedCounter: HTMLElement =
      (document.getElementById('images_accepted') as HTMLElement);
  private rejectionReasonCells: HTMLElement[] = [];
  private fieldImage: HTMLImageElement = new Image();

  constructor(private readonly connection: Connection) {
    (document.getElementById('field') as HTMLElement).appendChild(this.canvas);

    this.fieldImage.src = "2023.png";

    for (const value in RejectionReason) {
      // Typescript generates an iterator that produces both numbers and
      // strings... don't do anything on the string iterations.
      if (isNaN(Number(value))) {
        continue;
      }
      const row = document.createElement('div');
      const nameCell = document.createElement('div');
      nameCell.innerHTML = RejectionReason[value];
      row.appendChild(nameCell);
      const valueCell = document.createElement('div');
      valueCell.innerHTML = 'NA';
      this.rejectionReasonCells.push(valueCell);
      row.appendChild(valueCell);
      document.getElementById('vision_readouts').appendChild(row);
    }

    for (let ii = 0; ii < PI_COLORS.length; ++ii) {
      const legendEntry = document.createElement('div');
      legendEntry.style.color = PI_COLORS[ii];
      legendEntry.innerHTML = 'PI' + (ii + 1).toString()
      document.getElementById('legend').appendChild(legendEntry);
    }

    this.connection.addConfigHandler(() => {
      // Visualization message is reliable so that we can see *all* the vision
      // matches.
      for (const pi in PIS) {
        this.connection.addReliableHandler(
            '/' + PIS[pi] + '/camera', "y2023.localizer.Visualization",
            (data) => {
              this.handleLocalizerDebug(pi, data);
            });
      }
      this.connection.addHandler(
          '/drivetrain', "frc971.control_loops.drivetrain.Status", (data) => {
            this.handleDrivetrainStatus(data);
          });
      this.connection.addHandler(
               '/localizer', "frc971.controls.LocalizerOutput", (data) => {
            this.handleLocalizerOutput(data);
          });
    });
  }

  private handleLocalizerDebug(pi: string, data: Uint8Array): void {
    const now = Date.now() / 1000.0;

    const fbBuffer = new ByteBuffer(data);
    this.localizerImageMatches.set(
        now, Visualization.getRootAsVisualization(fbBuffer));

    const debug = this.localizerImageMatches.get(now);

    if (debug.statistics()) {
      if (debug.statistics().rejectionReasonsLength() ==
          this.rejectionReasonCells.length) {
        for (let ii = 0; ii < debug.statistics().rejectionReasonsLength();
             ++ii) {
          this.rejectionReasonCells[ii].innerHTML =
              debug.statistics().rejectionReasons(ii).count().toString();
        }
      } else {
        console.error('Unexpected number of rejection reasons in counter.');
      }
    }
  }

  private handleLocalizerOutput(data: Uint8Array): void {
    const fbBuffer = new ByteBuffer(data);
    this.localizerOutput = LocalizerOutput.getRootAsLocalizerOutput(fbBuffer);
  }

  private handleDrivetrainStatus(data: Uint8Array): void {
    const fbBuffer = new ByteBuffer(data);
    this.drivetrainStatus = DrivetrainStatus.getRootAsStatus(fbBuffer);
  }

  drawField(): void {
    const ctx = this.canvas.getContext('2d');
    ctx.save();
    ctx.scale(1.0, -1.0);
    ctx.drawImage(
        this.fieldImage, 0, 0, this.fieldImage.width, this.fieldImage.height,
        -FIELD_EDGE_X, -FIELD_SIDE_Y, FIELD_LENGTH, FIELD_WIDTH);
    ctx.restore();
  }

  drawCamera(
      x: number, y: number, theta: number, color: string = 'blue'): void {
    const ctx = this.canvas.getContext('2d');
    ctx.save();
    ctx.translate(x, y);
    ctx.rotate(theta);
    ctx.strokeStyle = color;
    ctx.beginPath();
    ctx.moveTo(0.5, 0.5);
    ctx.lineTo(0, 0);
    ctx.lineTo(0.5, -0.5);
    ctx.stroke();
    ctx.beginPath();
    ctx.arc(0, 0, 0.25, -Math.PI / 4, Math.PI / 4);
    ctx.stroke();
    ctx.restore();
  }

  drawRobot(
      x: number, y: number, theta: number,
      color: string = 'blue', dashed: boolean = false): void {
    const ctx = this.canvas.getContext('2d');
    ctx.save();
    ctx.translate(x, y);
    ctx.rotate(theta);
    ctx.strokeStyle = color;
    ctx.lineWidth = ROBOT_WIDTH / 10.0;
    if (dashed) {
      ctx.setLineDash([0.05, 0.05]);
    } else {
      // Empty array = solid line.
      ctx.setLineDash([]);
    }
    ctx.rect(-ROBOT_LENGTH / 2, -ROBOT_WIDTH / 2, ROBOT_LENGTH, ROBOT_WIDTH);
    ctx.stroke();

    // Draw line indicating which direction is forwards on the robot.
    ctx.beginPath();
    ctx.moveTo(0, 0);
    ctx.lineTo(ROBOT_LENGTH / 2.0, 0);
    ctx.stroke();

    ctx.restore();
  }

  setZeroing(div: HTMLElement): void {
    div.innerHTML = 'zeroing';
    div.classList.remove('faulted');
    div.classList.add('zeroing');
    div.classList.remove('near');
  }

  setValue(div: HTMLElement, val: number): void {
    div.innerHTML = val.toFixed(4);
    div.classList.remove('faulted');
    div.classList.remove('zeroing');
    div.classList.remove('near');
  }

  draw(): void {
    this.reset();
    this.drawField();

    // Draw the matches with debugging information from the localizer.
    const now = Date.now() / 1000.0;

    if (this.drivetrainStatus && this.drivetrainStatus.trajectoryLogging()) {
      this.drawRobot(
          this.drivetrainStatus.trajectoryLogging().x(),
          this.drivetrainStatus.trajectoryLogging().y(),
          this.drivetrainStatus.trajectoryLogging().theta(), "#000000FF",
          false);
    }

    if (this.localizerOutput) {
      if (!this.localizerOutput.zeroed()) {
        this.setZeroing(this.x);
        this.setZeroing(this.y);
        this.setZeroing(this.theta);
      } else {
        this.setValue(this.x, this.localizerOutput.x());
        this.setValue(this.y, this.localizerOutput.y());
        this.setValue(this.theta, this.localizerOutput.theta());
      }

      this.drawRobot(
          this.localizerOutput.x(), this.localizerOutput.y(),
          this.localizerOutput.theta());

      this.imagesAcceptedCounter.innerHTML =
          this.localizerOutput.imageAcceptedCount().toString();
    }

    for (const [time, value] of this.localizerImageMatches) {
      const age = now - time;
      const kRemovalAge = 1.0;
      if (age > kRemovalAge) {
        this.localizerImageMatches.delete(time);
        continue;
      }
      const kMaxImageAlpha = 0.5;
      const ageAlpha = kMaxImageAlpha * (kRemovalAge - age) / kRemovalAge
      for (let i = 0; i < value.targetsLength(); i++) {
        const imageDebug = value.targets(i);
        const x = imageDebug.impliedRobotX();
        const y = imageDebug.impliedRobotY();
        const theta = imageDebug.impliedRobotTheta();
        const cameraX = imageDebug.cameraX();
        const cameraY = imageDebug.cameraY();
        const cameraTheta = imageDebug.cameraTheta();
        const accepted = imageDebug.accepted();
        // Make camera readings fade over time.
        const alpha = Math.round(255 * ageAlpha).toString(16).padStart(2, '0');
        const dashed = false;
        const cameraRgb = PI_COLORS[imageDebug.camera()];
        const cameraRgba = cameraRgb + alpha;
        this.drawRobot(x, y, theta, cameraRgba, dashed);
        this.drawCamera(cameraX, cameraY, cameraTheta, cameraRgba);
      }
    }

    window.requestAnimationFrame(() => this.draw());
  }

  reset(): void {
    const ctx = this.canvas.getContext('2d');
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    const size = window.innerHeight * 0.9;
    ctx.canvas.height = size;
    const width = size / 2 + 20;
    ctx.canvas.width = width;
    ctx.clearRect(0, 0, size, width);

    // Translate to center of display.
    ctx.translate(width / 2, size / 2);
    // Coordinate system is:
    // x -> forward.
    // y -> to the left.
    ctx.rotate(-Math.PI / 2);
    ctx.scale(1, -1);

    const M_TO_PX = (size - 10) / FIELD_LENGTH;
    ctx.scale(M_TO_PX, M_TO_PX);
    ctx.lineWidth = 1 / M_TO_PX;
  }
}
