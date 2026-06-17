#include "handlers.h"

PlanHandler::PlanHandler(Database& database) : db(database) {}

void PlanHandler::registerRoutes(crow::SimpleApp& app) {
    CROW_ROUTE(app, "/plans").methods("POST"_method)
        ([this](const crow::request& req) { return createPlan(req); });

    CROW_ROUTE(app, "/plans").methods("GET"_method)
        ([this](const crow::request& req) { return getPlans(req); });

    CROW_ROUTE(app, "/plans/range").methods("GET"_method)
        ([this](const crow::request& req) { return getPlansByRange(req); });

    CROW_ROUTE(app, "/plans/<int>").methods("DELETE"_method)
        ([this](int id) { return deletePlan(id); });

    CROW_ROUTE(app, "/health")
        ([this]() { return healthCheck(); });
}

crow::response PlanHandler::createPlan(const crow::request& req) {
    auto body = crow::json::load(req.body);
    if (!body) {
        return crow::response(400, R"({"error": "Invalid JSON"})");
    }

    std::string title = body["title"].s();
    std::string description = body["description"].s();
    std::string date = body["date"].s();
    std::string time = body["time"].s();
    bool is_all_day = body["is_all_day"].b();

    if (title.empty() || date.empty()) {
        return crow::response(400, R"({"error": "title and date are required"})");
    }

    int id = db.createPlan(title, description, date, time, is_all_day);
    if (id == -1) {
        return crow::response(500, R"({"error": "Failed to create plan"})");
    }

    crow::json::wvalue result;
    result["id"] = id;
    result["message"] = "Plan created";

    return crow::response(201, result);
}

crow::response PlanHandler::getPlans(const crow::request& req) {
    auto date = req.url_params.get("date");

    std::vector<Plan> plans;
    if (date) {
        plans = db.getPlansByDate(date);
    } else {
        plans = db.getAllPlans();
    }

    return crow::response(200, plansToJson(plans));
}

crow::response PlanHandler::getPlansByRange(const crow::request& req) {
    auto start = req.url_params.get("start");
    auto end = req.url_params.get("end");

    if (!start || !end) {
        return crow::response(400, R"({"error": "start and end params required"})");
    }

    auto plans = db.getPlansByDateRange(start, end);
    return crow::response(200, plansToJson(plans));
}

crow::response PlanHandler::deletePlan(int id) {
    if (db.deletePlan(id)) {
        return crow::response(200, R"({"message": "Plan deleted"})");
    }
    return crow::response(404, R"({"error": "Plan not found"})");
}

crow::response PlanHandler::healthCheck() {
    crow::json::wvalue result;
    result["status"] = "ok";
    result["service"] = "todo-backend";
    return crow::response(200, result);
}

crow::json::wvalue PlanHandler::plansToJson(const std::vector<Plan>& plans) {
    crow::json::wvalue result = crow::json::wvalue::list();
    for (size_t i = 0; i < plans.size(); ++i) {
        result[i]["id"] = plans[i].id;
        result[i]["title"] = plans[i].title;
        result[i]["description"] = plans[i].description;
        result[i]["date"] = plans[i].date;
        result[i]["time"] = plans[i].time;
        result[i]["is_all_day"] = plans[i].is_all_day;
        result[i]["created_at"] = plans[i].created_at;
    }
    return result;
}
