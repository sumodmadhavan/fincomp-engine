const std = @import("std");
const math = std.math;
const print = std.debug.print;
const time = std.time;

const FinancialParams = struct {
    num_years: i32,
    au_hours: f64,
    initial_tsn: f64,
    rate_escalation: f64,
    aic: f64,
    hsi_tsn: f64,
    overhaul_tsn: f64,
    hsi_cost: f64,
    overhaul_cost: f64,
};

fn calculateFinancials(rate: f64, params: *const FinancialParams) f64 {
    var cumulative_profit: f64 = 0.0;
    var year: i32 = 1;
    while (year <= params.num_years) : (year += 1) {
        const tsn = params.initial_tsn + params.au_hours * @as(f64, @floatFromInt(year));
        const escalated_rate = rate * math.pow(f64, 1.0 + params.rate_escalation / 100.0, @as(f64, @floatFromInt(year - 1)));

        const engine_revenue = params.au_hours * escalated_rate;
        const total_revenue = engine_revenue * (1.0 + params.aic / 100.0);

        const hsi_cost = if (tsn >= params.hsi_tsn and (year == 1 or tsn - params.au_hours < params.hsi_tsn)) params.hsi_cost else 0.0;
        const overhaul_cost = if (tsn >= params.overhaul_tsn and (year == 1 or tsn - params.au_hours < params.overhaul_tsn)) params.overhaul_cost else 0.0;

        cumulative_profit += total_revenue - (hsi_cost + overhaul_cost);
    }
    return cumulative_profit;
}

fn objectiveFunction(rate: f64, params: *const FinancialParams) f64 {
    return calculateFinancials(rate, params) - 3000000.0; // Target profit
}

fn derivativeFunction(rate: f64, params: *const FinancialParams) f64 {
    const epsilon = 1e-6;
    return (objectiveFunction(rate + epsilon, params) - objectiveFunction(rate, params)) / epsilon;
}

fn newtonRaphson(
    comptime f: fn (f64, *const FinancialParams) f64,
    comptime df: fn (f64, *const FinancialParams) f64,
    x0: f64,
    tol: f64,
    max_iter: i32,
    params: *const FinancialParams,
) f64 {
    var x = x0;
    var i: i32 = 0;
    while (i < max_iter) : (i += 1) {
        const fx = f(x, params);
        if (math.fabs(fx) < tol) {
            return x;
        }
        const dfx = df(x, params);
        if (dfx == 0) {
            print("Derivative is zero. Newton-Raphson method failed.\n", .{});
            return x;
        }
        x -= fx / dfx;
    }
    print("Newton-Raphson method did not converge within {d} iterations\n", .{max_iter});
    return x;
}

pub fn main() !void {
    const params = FinancialParams{
        .num_years = 100,
        .au_hours = 450.0,
        .initial_tsn = 100.0,
        .rate_escalation = 5.0,
        .aic = 10.0,
        .hsi_tsn = 1000.0,
        .overhaul_tsn = 3000.0,
        .hsi_cost = 50000.0,
        .overhaul_cost = 100000.0,
    };

    const initial_rate = 100.0;
    const target_profit = 3000000.0;

    const start = time.nanoTimestamp();

    const initial_cumulative_profit = calculateFinancials(initial_rate, &params);
    print("Initial Warranty Rate: {d:.2}\n", .{initial_rate});
    print("Initial Cumulative Profit: {d:.2}\n", .{initial_cumulative_profit});

    const optimal_rate = newtonRaphson(objectiveFunction, derivativeFunction, initial_rate, 1e-8, 100, &params);
    print("\nOptimal Warranty Rate to achieve {d:.2} profit: {d:.7}\n", .{ target_profit, optimal_rate });

    const final_cumulative_profit = calculateFinancials(optimal_rate, &params);
    print("\nFinal Cumulative Profit: {d:.2}\n", .{final_cumulative_profit});

    const end = time.nanoTimestamp();
    const execution_time_ns = @as(u64, @intCast(end - start));
    const execution_time_us = execution_time_ns / 1000;
    print("\nExecution time: {d} microseconds\n", .{execution_time_us});
}
