#include <gtest/gtest.h>
#include <string>

#include "models.h"

TEST(ModelTest, PlanDefaultValues) {
    Plan plan;
    plan.id = 0;
    plan.title = "";
    plan.description = "";
    plan.date = "";
    plan.time = "";
    plan.is_all_day = false;
    plan.created_at = "";

    EXPECT_EQ(plan.id, 0);
    EXPECT_TRUE(plan.title.empty());
    EXPECT_TRUE(plan.description.empty());
    EXPECT_TRUE(plan.date.empty());
    EXPECT_TRUE(plan.time.empty());
    EXPECT_FALSE(plan.is_all_day);
    EXPECT_TRUE(plan.created_at.empty());
}

TEST(ModelTest, PlanWithValues) {
    Plan plan;
    plan.id = 1;
    plan.title = "Test Plan";
    plan.description = "Test Description";
    plan.date = "2024-12-25";
    plan.time = "14:00";
    plan.is_all_day = false;
    plan.created_at = "2024-12-24T10:00:00Z";

    EXPECT_EQ(plan.id, 1);
    EXPECT_EQ(plan.title, "Test Plan");
    EXPECT_EQ(plan.description, "Test Description");
    EXPECT_EQ(plan.date, "2024-12-25");
    EXPECT_EQ(plan.time, "14:00");
    EXPECT_FALSE(plan.is_all_day);
    EXPECT_EQ(plan.created_at, "2024-12-24T10:00:00Z");
}

TEST(ModelTest, PlanAllDay) {
    Plan plan;
    plan.id = 2;
    plan.title = "All Day";
    plan.description = "No time";
    plan.date = "2024-12-26";
    plan.time = "";
    plan.is_all_day = true;

    EXPECT_TRUE(plan.is_all_day);
    EXPECT_TRUE(plan.time.empty());
}

TEST(ModelTest, PlanRequestDefaultValues) {
    PlanRequest req;
    req.title = "";
    req.description = "";
    req.date = "";
    req.time = "";
    req.is_all_day = false;

    EXPECT_TRUE(req.title.empty());
    EXPECT_TRUE(req.description.empty());
    EXPECT_TRUE(req.date.empty());
    EXPECT_TRUE(req.time.empty());
    EXPECT_FALSE(req.is_all_day);
}

TEST(ModelTest, PlanRequestWithValues) {
    PlanRequest req;
    req.title = "New Plan";
    req.description = "New Description";
    req.date = "2024-12-25";
    req.time = "14:00";
    req.is_all_day = false;

    EXPECT_EQ(req.title, "New Plan");
    EXPECT_EQ(req.description, "New Description");
    EXPECT_EQ(req.date, "2024-12-25");
    EXPECT_EQ(req.time, "14:00");
    EXPECT_FALSE(req.is_all_day);
}

TEST(ModelTest, PlanRequestAllDay) {
    PlanRequest req;
    req.title = "All Day Event";
    req.description = "No specific time";
    req.date = "2024-12-26";
    req.time = "";
    req.is_all_day = true;

    EXPECT_TRUE(req.is_all_day);
    EXPECT_TRUE(req.time.empty());
}

TEST(ModelTest, PlanCopySemantics) {
    Plan plan1;
    plan1.id = 1;
    plan1.title = "Original";

    Plan plan2 = plan1;
    EXPECT_EQ(plan2.id, 1);
    EXPECT_EQ(plan2.title, "Original");

    plan2.title = "Modified";
    EXPECT_EQ(plan1.title, "Original");
    EXPECT_EQ(plan2.title, "Modified");
}

TEST(ModelTest, PlanRequestCopySemantics) {
    PlanRequest req1;
    req1.title = "Original";

    PlanRequest req2 = req1;
    EXPECT_EQ(req2.title, "Original");

    req2.title = "Modified";
    EXPECT_EQ(req1.title, "Original");
    EXPECT_EQ(req2.title, "Modified");
}

TEST(ModelTest, PlanSpecialCharacters) {
    Plan plan;
    plan.title = "Тест с кириллицей";
    plan.description = "Описание на русском";
    plan.date = "2024-12-25";
    plan.time = "14:00";

    EXPECT_EQ(plan.title, "Тест с кириллицей");
    EXPECT_EQ(plan.description, "Описание на русском");
}

TEST(ModelTest, PlanLongStrings) {
    Plan plan;
    plan.title = std::string(1000, 'A');
    plan.description = std::string(5000, 'B');

    EXPECT_EQ(plan.title.length(), 1000);
    EXPECT_EQ(plan.description.length(), 5000);
}
