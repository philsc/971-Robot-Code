#include <unistd.h>
#include <memory>
#include <random>

#include "aos/die.h"
#include "aos/flatbuffers.h"
#include "aos/json_to_flatbuffer.h"
#include "frc971/control_loops/position_sensor_sim.h"
#include "frc971/control_loops/team_number_test_environment.h"
#include "frc971/zeroing/zeroing.h"
#include "glog/logging.h"
#include "gtest/gtest.h"
#include "y2017/constants.h"
#include "y2017/control_loops/superstructure/column/column_zeroing.h"

namespace y2017 {
namespace control_loops {
namespace superstructure {
namespace column {

using ::frc971::HallEffectAndPosition;
using ::frc971::control_loops::PositionSensorSimulator;
using ::frc971::zeroing::HallEffectAndPositionZeroingEstimator;

class ZeroingTest : public ::testing::Test {
 public:
  explicit ZeroingTest()
      : indexer_sensor_(2.0 * M_PI),
        turret_sensor_(2.0 * M_PI),
        column_zeroing_estimator_(constants::GetValues().column) {}

 protected:
  void SetUp() override { aos::SetDieTestMode(true); }

  PositionSensorSimulator indexer_sensor_;
  PositionSensorSimulator turret_sensor_;
  ColumnZeroingEstimator column_zeroing_estimator_;

  void InitializeHallEffectAndPosition(double indexer_start_position,
                                       double turret_start_position) {
    indexer_sensor_.InitializeHallEffectAndPosition(
        indexer_start_position,
        constants::GetValues().column.indexer_zeroing.lower_hall_position,
        constants::GetValues().column.indexer_zeroing.upper_hall_position);
    turret_sensor_.InitializeHallEffectAndPosition(
        turret_start_position - indexer_start_position,
        constants::GetValues().column.turret_zeroing.lower_hall_position,
        constants::GetValues().column.turret_zeroing.upper_hall_position);
  }

  void MoveTo(double indexer, double turret) {
    flatbuffers::FlatBufferBuilder fbb;
    indexer_sensor_.MoveTo(indexer);
    turret_sensor_.MoveTo(turret - indexer);

    HallEffectAndPosition::Builder indexer_builder(fbb);
    flatbuffers::Offset<HallEffectAndPosition> indexer_offset =
        indexer_sensor_.GetSensorValues(&indexer_builder);

    HallEffectAndPosition::Builder turret_builder(fbb);
    flatbuffers::Offset<HallEffectAndPosition> turret_offset =
        turret_sensor_.GetSensorValues(&turret_builder);

    ColumnPosition::Builder column_position_builder(fbb);
    column_position_builder.add_indexer(indexer_offset);
    column_position_builder.add_turret(turret_offset);
    fbb.Finish(column_position_builder.Finish());


    aos::FlatbufferDetachedBuffer<ColumnPosition> column_position(fbb.Release());
    LOG(INFO) << "Position: " << aos::FlatbufferToJson(column_position);

    column_zeroing_estimator_.UpdateEstimate(column_position.message());
  }
};

// Tests that starting at 0 and spinning around the indexer in a full circle
// zeroes it (skipping the only offset_ready state).
TEST_F(ZeroingTest, TestZeroStartingPoint) {
  InitializeHallEffectAndPosition(0.0, 0.0);

  for (double i = 0; i < 2.0 * M_PI; i += M_PI / 100) {
    MoveTo(i, 0.0);
    EXPECT_EQ(column_zeroing_estimator_.zeroed(),
              column_zeroing_estimator_.offset_ready());
  }
  EXPECT_NEAR(column_zeroing_estimator_.indexer_offset(), 0.0, 1e-6);
  EXPECT_NEAR(column_zeroing_estimator_.turret_offset(), 0.0, 1e-6);
  EXPECT_TRUE(column_zeroing_estimator_.zeroed());
}

// Tests that starting the turret at M_PI + a bit stays offset_ready until we
// move the turret closer to 0.0
TEST_F(ZeroingTest, TestPiStartingPoint) {
  constexpr double kStartingPosition = M_PI + 0.1;
  InitializeHallEffectAndPosition(0.0, kStartingPosition);

  for (double i = 0; i < 2.0 * M_PI; i += M_PI / 100) {
    MoveTo(i, kStartingPosition);
    EXPECT_FALSE(column_zeroing_estimator_.zeroed());
  }
  for (double i = kStartingPosition; i > 0.0; i -= M_PI / 100) {
    MoveTo(2.0 * M_PI, i);
  }
  EXPECT_TRUE(column_zeroing_estimator_.offset_ready());
  EXPECT_TRUE(column_zeroing_estimator_.zeroed());

  EXPECT_NEAR(column_zeroing_estimator_.indexer_offset(), 0.0, 1e-6);
  EXPECT_NEAR(column_zeroing_estimator_.turret_offset(), kStartingPosition,
              1e-6);
  EXPECT_TRUE(column_zeroing_estimator_.zeroed());
}

// Tests that starting the turret at -M_PI - a bit stays offset_ready until we
// move the turret closer to 0.0
TEST_F(ZeroingTest, TestMinusPiStartingPoint) {
  constexpr double kStartingPosition = -M_PI - 0.1;
  InitializeHallEffectAndPosition(0.0, kStartingPosition);

  for (double i = 0; i < 2.0 * M_PI; i += M_PI / 100) {
    MoveTo(i, kStartingPosition);
    EXPECT_FALSE(column_zeroing_estimator_.zeroed());
  }
  for (double i = kStartingPosition; i < 0.0; i += M_PI / 100) {
    MoveTo(2.0 * M_PI, i);
  }
  EXPECT_TRUE(column_zeroing_estimator_.offset_ready());
  EXPECT_TRUE(column_zeroing_estimator_.zeroed());

  EXPECT_NEAR(column_zeroing_estimator_.indexer_offset(), 0.0, 1e-6);
  EXPECT_NEAR(column_zeroing_estimator_.turret_offset(), kStartingPosition,
              1e-6);
  EXPECT_TRUE(column_zeroing_estimator_.zeroed());
}

}  // namespace column
}  // namespace superstructure
}  // namespace control_loops
}  // namespace y2017
