#pragma once

#include <string>

struct Plan {
    int id;
    std::string title;
    std::string description;
    std::string date;
    std::string time;
    bool is_all_day;
    std::string created_at;
};

struct PlanRequest {
    std::string title;
    std::string description;
    std::string date;
    std::string time;
    bool is_all_day;
};
