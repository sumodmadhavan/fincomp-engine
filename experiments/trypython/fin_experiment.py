"""
This module [brief description].

[Detailed description]
"""

import time

import numpy as np
from scipy.optimize import root_scalar


def calculate_financials(
    rate,
    num_years=40,
    au_hours=450,
    initial_tsn=100,
    rate_escalation=5,
    aic=10,
    hsi_tsn=1000,
    overhaul_tsn=3000,
    hsi_cost=50000,
    overhaul_cost=100000,
):
    years = np.arange(1, num_years + 1)
    tsn = initial_tsn + au_hours * years

    escalated_rate = rate * (1 + rate_escalation / 100) ** (years - 1)
    engine_revenue = au_hours * escalated_rate
    aic_revenue = engine_revenue * aic / 100
    total_revenue = engine_revenue + aic_revenue

    hsi = np.zeros(num_years)
    hsi[np.where((tsn >= hsi_tsn) & (np.roll(tsn, 1) < hsi_tsn))[0]] = 1
    hsi[0] = 1 if tsn[0] >= hsi_tsn else 0

    overhaul = np.zeros(num_years)
    overhaul[np.where((tsn >= overhaul_tsn) & (np.roll(tsn, 1) < overhaul_tsn))[0]] = 1
    overhaul[0] = 1 if tsn[0] >= overhaul_tsn else 0

    total_cost = hsi * hsi_cost + overhaul * overhaul_cost
    total_profit = total_revenue - total_cost
    cumulative_profit = np.sum(total_profit)

    return cumulative_profit


def goal_seek(target_profit, initial_guess=320):
    def objective(rate):
        return calculate_financials(rate) - target_profit

    result = root_scalar(
        objective, x0=initial_guess, x1=initial_guess * 1.1, method="secant", xtol=1e-8
    )

    return result.root, result.iterations


def main():
    target_profit = 3_000_000
    initial_rate = 320

    start_time = time.time()

    optimal_rate, iterations = goal_seek(target_profit, initial_rate)

    end_time = time.time()
    execution_time = end_time - start_time

    print(f"Optimal rate: {optimal_rate:.7f}")
    print(f"Final profit: {calculate_financials(optimal_rate):.2f}")
    print(f"\nInitial rate: {initial_rate}")
    print(f"Initial profit: {calculate_financials(initial_rate):.2f}")

    print(f"\nExecution time: {execution_time:.6f} seconds")
    print(f"Number of iterations: {iterations}")


if __name__ == "__main__":
    main()
