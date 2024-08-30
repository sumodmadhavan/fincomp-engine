import numpy as np
from numba import njit
from scipy.optimize import brentq
import time

@njit
def calculate_financials(rate, num_years=10, au_hours=450, initial_tsn=100,
                         rate_escalation=5, aic=10, hsi_tsn=1000, overhaul_tsn=3000,
                         hsi_cost=50000, overhaul_cost=100000):
    years = np.arange(1, num_years + 1)
    tsn = initial_tsn + au_hours * years
    
    escalated_rate = rate * (1 + rate_escalation/100) ** (years - 1)
    engine_revenue = au_hours * escalated_rate
    aic_revenue = engine_revenue * aic / 100
    total_revenue = engine_revenue + aic_revenue
    
    hsi = np.zeros(num_years)
    for i in range(num_years):
        if i == 0:
            hsi[i] = 1 if tsn[i] >= hsi_tsn else 0
        else:
            hsi[i] = 1 if tsn[i] >= hsi_tsn and tsn[i-1] < hsi_tsn else 0
    
    overhaul = np.zeros(num_years)
    for i in range(num_years):
        if i == 0:
            overhaul[i] = 1 if tsn[i] >= overhaul_tsn else 0
        else:
            overhaul[i] = 1 if tsn[i] >= overhaul_tsn and tsn[i-1] < overhaul_tsn else 0
    
    total_cost = hsi * hsi_cost + overhaul * overhaul_cost
    total_profit = total_revenue - total_cost
    cumulative_profit = np.cumsum(total_profit)
    
    return cumulative_profit[-1]

@njit
def objective(rate, target_profit):
    return calculate_financials(rate) - target_profit

def goal_seek(target_profit, lower_bound=1, upper_bound=1000):
    return brentq(objective, lower_bound, upper_bound, args=(target_profit,), xtol=1e-8)

def main():
    target_profit = 3_000_000
    initial_rate = 320

    start_time = time.time()
    
    # Warm-up JIT compilation
    calculate_financials(initial_rate)
    objective(initial_rate, target_profit)
    
    optimal_rate = goal_seek(target_profit)
    
    end_time = time.time()
    execution_time = end_time - start_time
    
    print(f"Optimal rate: {optimal_rate:.7f}")
    print(f"Final profit: {calculate_financials(optimal_rate):.2f}")
    print(f"\nInitial rate: {initial_rate}")
    print(f"Initial profit: {calculate_financials(initial_rate):.2f}")
    
    print(f"\nExecution time: {execution_time:.6f} seconds")

if __name__ == "__main__":
    main()
