#ifndef Y2023_CONTROL_LOOPS_DRIVETRAIN_DRIVETRAIN_BASE_H_
#define Y2023_CONTROL_LOOPS_DRIVETRAIN_DRIVETRAIN_BASE_H_

#include "frc971/control_loops/drivetrain/drivetrain_config.h"

namespace y2023 {
namespace control_loops {
namespace drivetrain {

const ::frc971::control_loops::drivetrain::DrivetrainConfig<double>
    &GetDrivetrainConfig();

}  // namespace drivetrain
}  // namespace control_loops
}  // namespace y2023

#endif  // Y2023_CONTROL_LOOPS_DRIVETRAIN_DRIVETRAIN_BASE_H_
