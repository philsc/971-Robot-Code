#ifndef Y2022_BOT3_CONTROL_LOOPS_SUPERSTRUCTURE_SUPERSTRUCTURE_H_
#define Y2022_BOT3_CONTROL_LOOPS_SUPERSTRUCTURE_SUPERSTRUCTURE_H_

#include "aos/events/event_loop.h"
#include "frc971/control_loops/control_loop.h"
#include "frc971/control_loops/drivetrain/drivetrain_status_generated.h"
#include "y2022_bot3/constants.h"
#include "y2022_bot3/control_loops/superstructure/superstructure_goal_generated.h"
#include "y2022_bot3/control_loops/superstructure/superstructure_output_generated.h"
#include "y2022_bot3/control_loops/superstructure/superstructure_position_generated.h"
#include "y2022_bot3/control_loops/superstructure/superstructure_status_generated.h"

namespace y2022_bot3 {
namespace control_loops {
namespace superstructure {

class Superstructure
    : public ::frc971::controls::ControlLoop<Goal, Position, Status, Output> {
 public:
  using RelativeEncoderSubsystem =
      ::frc971::control_loops::StaticZeroingSingleDOFProfiledSubsystem<
          ::frc971::zeroing::RelativeEncoderZeroingEstimator,
          ::frc971::control_loops::RelativeEncoderProfiledJointStatus>;

  using PotAndAbsoluteEncoderSubsystem =
      ::frc971::control_loops::StaticZeroingSingleDOFProfiledSubsystem<
          ::frc971::zeroing::PotAndAbsoluteEncoderZeroingEstimator,
          ::frc971::control_loops::PotAndAbsoluteEncoderProfiledJointStatus>;

  explicit Superstructure(::aos::EventLoop *event_loop,
                          std::shared_ptr<const constants::Values> values,
                          const ::std::string &name = "/superstructure");

  inline const PotAndAbsoluteEncoderSubsystem &intake() const {
    return intake_;
  }
  inline const PotAndAbsoluteEncoderSubsystem &climber_left() const {
    return climber_left_;
  }
  inline const PotAndAbsoluteEncoderSubsystem &climber_right() const {
    return climber_right_;
  }

 protected:
  virtual void RunIteration(const Goal *unsafe_goal, const Position *position,
                            aos::Sender<Output>::Builder *output,
                            aos::Sender<Status>::Builder *status) override;

 private:
  std::shared_ptr<const constants::Values> values_;

  PotAndAbsoluteEncoderSubsystem climber_left_;
  PotAndAbsoluteEncoderSubsystem climber_right_;
  PotAndAbsoluteEncoderSubsystem intake_;

  aos::Fetcher<frc971::control_loops::drivetrain::Status>
      drivetrain_status_fetcher_;
  aos::Fetcher<aos::JoystickState> joystick_state_fetcher_;

  aos::Alliance alliance_ = aos::Alliance::kInvalid;

  DISALLOW_COPY_AND_ASSIGN(Superstructure);
};

}  // namespace superstructure
}  // namespace control_loops
}  // namespace y2022_bot3

#endif  // Y2022_BOT3_CONTROL_LOOPS_SUPERSTRUCTURE_SUPERSTRUCTURE_H_
