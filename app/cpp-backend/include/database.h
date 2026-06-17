#pragma once

#include <string>
#include <vector>
#include <sqlite3.h>
#include "models.h"

class Database {
public:
    Database(const std::string& dbPath);
    ~Database();

    Database(const Database&) = delete;
    Database& operator=(const Database&) = delete;

    int createPlan(const std::string& title, const std::string& description,
                   const std::string& date, const std::string& time, bool isAllDay);

    std::vector<Plan> getPlansByDate(const std::string& date);
    std::vector<Plan> getPlansByDateRange(const std::string& startDate, const std::string& endDate);
    std::vector<Plan> getAllPlans();
    bool deletePlan(int id);

private:
    sqlite3* db;
    void createTable();
};
