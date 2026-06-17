#include <crow.h>
#include "database.h"
#include "handlers.h"

int main() {
    try {
        Database db("todo.db");
        crow::SimpleApp app;

        PlanHandler handler(db);
        handler.registerRoutes(app);

        app.port(8081).multithreaded().run();
    } catch (const std::exception& e) {
        return 1;
    }
    return 0;
}
