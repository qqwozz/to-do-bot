#pragma once

#include <crow.h>
#include "database.h"

class PlanHandler {
public:
    explicit PlanHandler(Database& database);
    void registerRoutes(crow::SimpleApp& app);

private:
    Database& db;

    crow::response createPlan(const crow::request& req);
    crow::response getPlans(const crow::request& req);
    crow::response getPlansByRange(const crow::request& req);
    crow::response deletePlan(int id);
    crow::response healthCheck();
    crow::json::wvalue plansToJson(const std::vector<Plan>& plans);
};
