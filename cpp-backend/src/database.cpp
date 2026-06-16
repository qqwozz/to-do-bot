#include "database.h"
#include <stdexcept>

Database::Database(const std::string& dbPath) {
    int rc = sqlite3_open(dbPath.c_str(), &db);
    if (rc) {
        throw std::runtime_error("Failed to open DB: " + std::string(sqlite3_errmsg(db)));
    }
    sqlite3_exec(db, "PRAGMA journal_mode=WAL;", nullptr, nullptr, nullptr);
    createTable();
}

Database::~Database() {
    if (db) sqlite3_close(db);
}

void Database::createTable() {
    const char* sql = R"(
        CREATE TABLE IF NOT EXISTS plans (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            date TEXT NOT NULL,
            time TEXT,
            is_all_day INTEGER DEFAULT 1,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_plans_date ON plans(date);
    )";

    char* errMsg = nullptr;
    int rc = sqlite3_exec(db, sql, nullptr, nullptr, &errMsg);
    if (rc != SQLITE_OK) {
        std::string error = errMsg;
        sqlite3_free(errMsg);
        throw std::runtime_error("Table creation failed: " + error);
    }
}

int Database::createPlan(const std::string& title, const std::string& description,
                         const std::string& date, const std::string& time, bool isAllDay) {
    const char* sql = "INSERT INTO plans (title, description, date, time, is_all_day) "
                      "VALUES (?, ?, ?, ?, ?)";

    sqlite3_stmt* stmt;
    int rc = sqlite3_prepare_v2(db, sql, -1, &stmt, nullptr);
    if (rc != SQLITE_OK) return -1;

    sqlite3_bind_text(stmt, 1, title.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, description.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 3, date.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 4, time.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_int(stmt, 5, isAllDay ? 1 : 0);

    rc = sqlite3_step(stmt);
    sqlite3_finalize(stmt);

    return (rc == SQLITE_DONE) ? static_cast<int>(sqlite3_last_insert_row_id(db)) : -1;
}

std::vector<Plan> Database::getPlansByDate(const std::string& date) {
    std::vector<Plan> plans;
    const char* sql = "SELECT * FROM plans WHERE date = ? ORDER BY is_all_day ASC, time ASC";

    sqlite3_stmt* stmt;
    int rc = sqlite3_prepare_v2(db, sql, -1, &stmt, nullptr);
    if (rc != SQLITE_OK) return plans;

    sqlite3_bind_text(stmt, 1, date.c_str(), -1, SQLITE_STATIC);

    while (sqlite3_step(stmt) == SQLITE_ROW) {
        Plan plan;
        plan.id = sqlite3_column_int(stmt, 0);
        plan.title = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 1));
        plan.description = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 2));
        plan.date = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 3));
        const char* time = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 4));
        plan.time = time ? time : "";
        plan.is_all_day = sqlite3_column_int(stmt, 5) == 1;
        const char* created = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 6));
        plan.created_at = created ? created : "";
        plans.push_back(plan);
    }

    sqlite3_finalize(stmt);
    return plans;
}

std::vector<Plan> Database::getPlansByDateRange(const std::string& startDate,
                                                const std::string& endDate) {
    std::vector<Plan> plans;
    const char* sql = "SELECT * FROM plans WHERE date >= ? AND date <= ? "
                      "ORDER BY date ASC, is_all_day ASC, time ASC";

    sqlite3_stmt* stmt;
    int rc = sqlite3_prepare_v2(db, sql, -1, &stmt, nullptr);
    if (rc != SQLITE_OK) return plans;

    sqlite3_bind_text(stmt, 1, startDate.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, endDate.c_str(), -1, SQLITE_STATIC);

    while (sqlite3_step(stmt) == SQLITE_ROW) {
        Plan plan;
        plan.id = sqlite3_column_int(stmt, 0);
        plan.title = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 1));
        plan.description = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 2));
        plan.date = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 3));
        const char* time = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 4));
        plan.time = time ? time : "";
        plan.is_all_day = sqlite3_column_int(stmt, 5) == 1;
        const char* created = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 6));
        plan.created_at = created ? created : "";
        plans.push_back(plan);
    }

    sqlite3_finalize(stmt);
    return plans;
}

std::vector<Plan> Database::getAllPlans() {
    std::vector<Plan> plans;
    const char* sql = "SELECT * FROM plans ORDER BY date ASC, is_all_day ASC, time ASC";

    sqlite3_stmt* stmt;
    int rc = sqlite3_prepare_v2(db, sql, -1, &stmt, nullptr);
    if (rc != SQLITE_OK) return plans;

    while (sqlite3_step(stmt) == SQLITE_ROW) {
        Plan plan;
        plan.id = sqlite3_column_int(stmt, 0);
        plan.title = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 1));
        plan.description = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 2));
        plan.date = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 3));
        const char* time = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 4));
        plan.time = time ? time : "";
        plan.is_all_day = sqlite3_column_int(stmt, 5) == 1;
        const char* created = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 6));
        plan.created_at = created ? created : "";
        plans.push_back(plan);
    }

    sqlite3_finalize(stmt);
    return plans;
}

bool Database::deletePlan(int id) {
    const char* sql = "DELETE FROM plans WHERE id = ?";

    sqlite3_stmt* stmt;
    int rc = sqlite3_prepare_v2(db, sql, -1, &stmt, nullptr);
    if (rc != SQLITE_OK) return false;

    sqlite3_bind_int(stmt, 1, id);
    rc = sqlite3_step(stmt);
    sqlite3_finalize(stmt);

    return (rc == SQLITE_DONE) && (sqlite3_changes(db) > 0);
}
