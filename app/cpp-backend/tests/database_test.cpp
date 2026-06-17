#include <gtest/gtest.h>
#include <sqlite3.h>
#include <string>
#include <filesystem>

#include "database.h"

class DatabaseTest : public ::testing::Test {
protected:
    std::string dbPath = "test_todo.db";
    Database* db;

    void SetUp() override {
        std::filesystem::remove(dbPath);
        db = new Database(dbPath);
    }

    void TearDown() override {
        delete db;
        std::filesystem::remove(dbPath);
    }
};

TEST_F(DatabaseTest, CreatePlan) {
    int id = db->createPlan("Встреча", "Обсудить проект", "2024-12-25", "14:00", false);
    EXPECT_GT(id, 0);
}

TEST_F(DatabaseTest, CreatePlanAllDay) {
    int id = db->createPlan("Дедлайн", "Сдать отчёт", "2024-12-26", "", true);
    EXPECT_GT(id, 0);

    auto plans = db->getPlansByDate("2024-12-26");
    ASSERT_EQ(plans.size(), 1);
    EXPECT_TRUE(plans[0].is_all_day);
    EXPECT_EQ(plans[0].title, "Дедлайн");
}

TEST_F(DatabaseTest, GetPlansByDate) {
    db->createPlan("План 1", "Описание 1", "2024-12-25", "10:00", false);
    db->createPlan("План 2", "Описание 2", "2024-12-25", "14:00", false);
    db->createPlan("План 3", "Описание 3", "2024-12-26", "10:00", false);

    auto plans = db->getPlansByDate("2024-12-25");
    EXPECT_EQ(plans.size(), 2);
}

TEST_F(DatabaseTest, GetPlansByDateSorted) {
    db->createPlan("Вечер", "Описание", "2024-12-25", "18:00", false);
    db->createPlan("Утро", "Описание", "2024-12-25", "08:00", false);
    db->createPlan("Весь день", "Описание", "2024-12-25", "", true);

    auto plans = db->getPlansByDate("2024-12-25");
    ASSERT_EQ(plans.size(), 3);

    EXPECT_TRUE(plans[0].is_all_day == false);
    EXPECT_EQ(plans[0].time, "08:00");

    EXPECT_TRUE(plans[2].is_all_day == true);
}

TEST_F(DatabaseTest, GetPlansByDateRange) {
    db->createPlan("План 1", "Описание", "2024-12-23", "10:00", false);
    db->createPlan("План 2", "Описание", "2024-12-25", "10:00", false);
    db->createPlan("План 3", "Описание", "2024-12-27", "10:00", false);
    db->createPlan("План 4", "Описание", "2024-12-30", "10:00", false);

    auto plans = db->getPlansByDateRange("2024-12-23", "2024-12-27");
    EXPECT_EQ(plans.size(), 3);
}

TEST_F(DatabaseTest, GetAllPlans) {
    db->createPlan("План 1", "Описание", "2024-12-25", "10:00", false);
    db->createPlan("План 2", "Описание", "2024-12-26", "10:00", false);

    auto plans = db->getAllPlans();
    EXPECT_EQ(plans.size(), 2);
}

TEST_F(DatabaseTest, DeletePlan) {
    int id = db->createPlan("Для удаления", "Описание", "2024-12-25", "10:00", false);
    EXPECT_TRUE(db->deletePlan(id));

    auto plans = db->getPlansByDate("2024-12-25");
    EXPECT_EQ(plans.size(), 0);
}

TEST_F(DatabaseTest, DeletePlanNotFound) {
    EXPECT_FALSE(db->deletePlan(9999));
}

TEST_F(DatabaseTest, EmptyDateReturnsEmpty) {
    auto plans = db->getPlansByDate("2099-01-01");
    EXPECT_EQ(plans.size(), 0);
}

TEST_F(DatabaseTest, PlanFields) {
    db->createPlan("Название", "Описание", "2024-12-25", "14:30", false);

    auto plans = db->getPlansByDate("2024-12-25");
    ASSERT_EQ(plans.size(), 1);

    EXPECT_EQ(plans[0].title, "Название");
    EXPECT_EQ(plans[0].description, "Описание");
    EXPECT_EQ(plans[0].date, "2024-12-25");
    EXPECT_EQ(plans[0].time, "14:30");
    EXPECT_FALSE(plans[0].is_all_day);
    EXPECT_FALSE(plans[0].created_at.empty());
}

TEST_F(DatabaseTest, MultipleInserts) {
    for (int i = 0; i < 10; i++) {
        int id = db->createPlan(
            "Plan " + std::to_string(i),
            "Desc",
            "2024-12-25",
            "10:00",
            false
        );
        EXPECT_GT(id, 0);
    }

    auto plans = db->getAllPlans();
    EXPECT_EQ(plans.size(), 10);
}
