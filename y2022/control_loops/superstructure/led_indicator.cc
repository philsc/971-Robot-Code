#include "y2022/control_loops/superstructure/led_indicator.h"

namespace led = ctre::phoenix::led;

namespace y2022::control_loops::superstructure {

LedIndicator::LedIndicator(aos::EventLoop *event_loop)
    : event_loop_(event_loop),
      drivetrain_output_fetcher_(
          event_loop_->MakeFetcher<frc971::control_loops::drivetrain::Output>(
              "/drivetrain")),
      superstructure_status_fetcher_(
          event_loop_->MakeFetcher<Status>("/superstructure")),
      server_statistics_fetcher_(
          event_loop_->MakeFetcher<aos::message_bridge::ServerStatistics>(
              "/roborio/aos")),
      client_statistics_fetcher_(
          event_loop_->MakeFetcher<aos::message_bridge::ClientStatistics>(
              "/roborio/aos")),
      localizer_output_fetcher_(
          event_loop_->MakeFetcher<frc971::controls::LocalizerOutput>(
              "/localizer")),
      gyro_reading_fetcher_(
          event_loop_->MakeFetcher<frc971::sensors::GyroReading>(
              "/drivetrain")) {
  led::CANdleConfiguration config;
  config.statusLedOffWhenActive = true;
  config.disableWhenLOS = false;
  config.brightnessScalar = 1.0;
  candle_.ConfigAllSettings(config, 0);

  event_loop_->AddPhasedLoop([&](int) { DecideColor(); },
                            std::chrono::milliseconds(20));
}

// This method will be called once per scheduler run
void LedIndicator::DisplayLed(uint8_t r, uint8_t g, uint8_t b) {
  candle_.SetLEDs(static_cast<int>(r), static_cast<int>(g),
                  static_cast<int>(b));
}

namespace {
bool DisconnectedPiServer(
    const aos::message_bridge::ServerStatistics &server_stats) {
  for (const auto *pi_server_status : *server_stats.connections()) {
    if (pi_server_status->state() == aos::message_bridge::State::DISCONNECTED &&
        pi_server_status->node()->name()->string_view() != "logger") {
      return true;
    }
  }
  return false;
}

bool DisconnectedPiClient(
    const aos::message_bridge::ClientStatistics &client_stats) {
  for (const auto *pi_client_status : *client_stats.connections()) {
    if (pi_client_status->state() == aos::message_bridge::State::DISCONNECTED &&
        pi_client_status->node()->name()->string_view() != "logger") {
      return true;
    }
  }
  return false;
}
}  // namespace

void LedIndicator::DecideColor() {
  superstructure_status_fetcher_.Fetch();
  server_statistics_fetcher_.Fetch();
  drivetrain_output_fetcher_.Fetch();
  client_statistics_fetcher_.Fetch();
  gyro_reading_fetcher_.Fetch();
  localizer_output_fetcher_.Fetch();

  if (localizer_output_fetcher_.get()) {
    if (localizer_output_fetcher_->image_accepted_count() !=
        last_accepted_count_) {
      last_accepted_count_ = localizer_output_fetcher_->image_accepted_count();
      last_accepted_time_ = event_loop_->monotonic_now();
    }
  }

  // Estopped
  if (superstructure_status_fetcher_.get() &&
      superstructure_status_fetcher_->estopped()) {
    DisplayLed(255, 0, 0);
    return;
  }

  // Not zeroed
  if (superstructure_status_fetcher_.get() &&
      !superstructure_status_fetcher_->zeroed()) {
    DisplayLed(255, 255, 0);
    return;
  }

  // If the imu gyro readings are not being sent/updated recently
  if (!gyro_reading_fetcher_.get() ||
      gyro_reading_fetcher_.context().monotonic_event_time <
          event_loop_->monotonic_now() -
              frc971::controls::kLoopFrequency * 10) {
    if (imu_flash_) {
      DisplayLed(255, 0, 0);
    } else {
      DisplayLed(255, 255, 255);
    }

    if (imu_counter_ % kFlashIterations == 0) {
      imu_flash_ = !imu_flash_;
    }
    imu_counter_++;
    return;
  }

  // Pi disconnected
  if ((server_statistics_fetcher_.get() &&
       DisconnectedPiServer(*server_statistics_fetcher_)) ||
      (client_statistics_fetcher_.get() &&
       DisconnectedPiClient(*client_statistics_fetcher_))) {
    if (disconnected_flash_) {
      DisplayLed(255, 0, 0);
    } else {
      DisplayLed(0, 255, 0);
    }

    if (disconnected_counter_ % kFlashIterations == 0) {
      disconnected_flash_ = !disconnected_flash_;
    }
    disconnected_counter_++;
    return;
  }

  // Statemachine
  if (superstructure_status_fetcher_.get()) {
    switch (superstructure_status_fetcher_->state()) {
      case (SuperstructureState::IDLE):
        DisplayLed(0, 0, 0);
        break;
      case (SuperstructureState::TRANSFERRING):
        DisplayLed(0, 0, 255);
        break;
      case (SuperstructureState::LOADING):
        DisplayLed(255, 255, 255);
        break;
      case (SuperstructureState::LOADED):
        if (!superstructure_status_fetcher_->ready_to_fire()) {
          DisplayLed(255, 140, 0);
        } else if (superstructure_status_fetcher_->front_intake_has_ball() ||
                   superstructure_status_fetcher_->back_intake_has_ball()) {
          DisplayLed(165, 42, 42);
        }
        break;
      case (SuperstructureState::SHOOTING):
        break;
    }

    if (event_loop_->monotonic_now() <
        last_accepted_time_ + std::chrono::seconds(2)) {
      DisplayLed(255, 0, 255);
    }
    return;
  }
}

}  // namespace y2022::control_loops::superstructure
