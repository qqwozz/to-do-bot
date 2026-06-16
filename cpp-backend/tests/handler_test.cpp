#include <gtest/gtest.h>
#include <crow.h>
#include <crow/http_request.h>
#include <string>
#include <memory>

#include "database.h"
#include "handlers.h"

class HandlerTest : public ::testing::Test {
protected:
    std::string dbPath = "test_handler.db";
    std::unique_ptr<Database> db;
    crow::SimpleApp app;

    void SetUp() override {
        std::filesystem::remove(dbPath);
        db = std::make_unique<Database>(dbPath);
    }

    void TearDown() override {
        db.reset();
        std::filesystem::remove(dbPath);
    }
};

TEST_F(HandlerTest, HealthCheck) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    auto resp = app.handle_request("GET", "/health", crow::request{});

    EXPECT_EQ(resp.code, 200);
}

TEST_F(HandlerTest, CreatePlanValid) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::json::wvalue body;
    body["title"] = "Test Plan";
    body["description"] = "Test Description";
    body["date"] = "2024-12-25";
    body["time"] = "14:00";
    body["is_all_day"] = false;

    crow::request req;
    req.body = crow::json::dump(body);
    req.set_header("Content-Type", "application/json");

    auto resp = app.handle_request("POST", "/plans", req);

    EXPECT_EQ(resp.code, 201);
}

TEST_F(HandlerTest, CreatePlanMissingTitle) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::json::wvalue body;
    body["description"] = "Test";
    body["date"] = "2024-12-25";

    crow::request req;
    req.body = crow::json::dump(body);
    req.set_header("Content-Type", "application/json");

    auto resp = app.handle_request("POST", "/plans", req);

    EXPECT_EQ(resp.code, 400);
}

TEST_F(HandlerTest, CreatePlanMissingDate) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::json::wvalue body;
    body["title"] = "Test";
    body["description"] = "Test";

    crow::request req;
    req.body = crow::json::dump(body);
    req.set_header("Content-Type", "application/json");

    auto resp = app.handle_request("POST", "/plans", req);

    EXPECT_EQ(resp.code, 400);
}

TEST_F(HandlerTest, CreatePlanInvalidJSON) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::request req;
    req.body = "invalid json";
    req.set_header("Content-Type", "application/json");

    auto resp = app.handle_request("POST", "/plans", req);

    EXPECT_EQ(resp.code, 400);
}

TEST_F(HandlerTest, GetPlansByDate) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    db->createPlan("Plan 1", "Desc 1", "2024-12-25", "10:00", false);
    db->createPlan("Plan 2", "Desc 2", "2024-12-25", "14:00", false);
    db->createPlan("Plan 3", "Desc 3", "2024-12-26", "10:00", false);

    crow::request req;
    req.url = "/plans?date=2024-12-25";

    auto resp = app.handle_request("GET", "/plans?date=2024-12-25", req);

    EXPECT_EQ(resp.code, 200);
}

TEST_F(HandlerTest, GetPlansAll) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    db->createPlan("Plan 1", "Desc 1", "2024-12-25", "10:00", false);
    db->createPlan("Plan 2", "Desc 2", "2024-12-26", "14:00", false);

    crow::request req;

    auto resp = app.handle_request("GET", "/plans", req);

    EXPECT_EQ(resp.code, 200);
}

TEST_F(HandlerTest, GetPlansByRange) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    db->createPlan("Plan 1", "Desc 1", "2024-12-23", "10:00", false);
    db->createPlan("Plan 2", "Desc 2", "2024-12-25", "10:00", false);
    db->createPlan("Plan 3", "Desc 3", "2024-12-27", "10:00", false);

    crow::request req;

    auto resp = app.handle_request("GET", "/plans/range?start=2024-12-23&end=2024-12-27", req);

    EXPECT_EQ(resp.code, 200);
}

TEST_F(HandlerTest, GetPlansByRangeMissingParams) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::request req;

    auto resp = app.handle_request("GET", "/plans/range", req);

    EXPECT_EQ(resp.code, 400);
}

TEST_F(HandlerTest, DeletePlan) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    int id = db->createPlan("To Delete", "Desc", "2024-12-25", "10:00", false);

    crow::request req;

    auto resp = app.handle_request("DELETE", "/plans/" + std::to_string(id), req);

    EXPECT_EQ(resp.code, 200);
}

TEST_F(HandlerTest, DeletePlanNotFound) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::request req;

    auto resp = app.handle_request("DELETE", "/plans/99999", req);

    EXPECT_EQ(resp.code, 404);
}

TEST_F(HandlerTest, CreateAndRetrievePlan) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::json::wvalue body;
    body["title"] = "Integration Test";
    body["description"] = "Test Description";
    body["date"] = "2024-12-25";
    body["time"] = "14:00";
    body["is_all_day"] = false;

    crow::request req;
    req.body = crow::json::dump(body);
    req.set_header("Content-Type", "application/json");

    auto createResp = app.handle_request("POST", "/plans", req);
    EXPECT_EQ(createResp.code, 201);

    auto plans = db->getPlansByDate("2024-12-25");
    EXPECT_EQ(plans.size(), 1);
    EXPECT_EQ(plans[0].title, "Integration Test");
}

TEST_F(HandlerTest, CreatePlanAllDay) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::json::wvalue body;
    body["title"] = "All Day Event";
    body["description"] = "No time";
    body["date"] = "2024-12-25";
    body["time"] = "";
    body["is_all_day"] = true;

    crow::request req;
    req.body = crow::json::dump(body);
    req.set_header("Content-Type", "application/json");

    auto resp = app.handle_request("POST", "/plans", req);
    EXPECT_EQ(resp.code, 201);

    auto plans = db->getPlansByDate("2024-12-25");
    ASSERT_EQ(plans.size(), 1);
    EXPECT_TRUE(plans[0].is_all_day);
}

TEST_F(HandlerTest, DeletePlanThenNotFound) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    int id = db->createPlan("To Delete", "Desc", "2024-12-25", "10:00", false);

    crow::request req;

    auto deleteResp = app.handle_request("DELETE", "/plans/" + std::to_string(id), req);
    EXPECT_EQ(deleteResp.code, 200);

    auto notFoundResp = app.handle_request("DELETE", "/plans/" + std::to_string(id), req);
    EXPECT_EQ(notFoundResp.code, 404);
}

TEST_F(HandlerTest, GetPlansByRangeEmpty) {
    PlanHandler handler(*db);
    handler.registerRoutes(app);

    crow::request req;

    auto resp = app.handle_request("GET", "/plans/range?start=2099-01-01&end=2099-01-31", req);

    EXPECT_EQ(resp.code, 200);
}
