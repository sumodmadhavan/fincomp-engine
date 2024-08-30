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
    total_revenue = au_hours * escalated_rate * (1 + aic/100)
    
    hsi = (tsn >= hsi_tsn) & (np.roll(tsn, 1) < hsi_tsn)
    hsi[0] = tsn[0] >= hsi_tsn
    overhaul = (tsn >= overhaul_tsn) & (np.roll(tsn, 1) < overhaul_tsn)
    overhaul[0] = tsn[0] >= overhaul_tsn
    
    total_cost = hsi * hsi_cost + overhaul * overhaul_cost
    return np.sum(total_revenue - total_cost)

@njit
def objective(rate, target_profit):
    return calculate_financials(rate) - target_profit

def goal_seek(target_profit, lower_bound=1, upper_bound=1000):
    return brentq(objective, lower_bound, upper_bound, args=(target_profit,), xtol=1e-8)

def main():
    target_profit = 3_000_000
    initial_rate = 320

    # Warm-up JIT compilation
    print("Warming up JIT compilation...")
    warm_up_start = time.time()
    calculate_financials(initial_rate)
    objective(initial_rate, target_profit)
    warm_up_end = time.time()
    print(f"Warm-up time: {warm_up_end - warm_up_start:.6f} seconds")

    # Actual timed run
    print("\nRunning timed calculation...")
    start_time = time.time()
    
    optimal_rate = goal_seek(target_profit)
    
    end_time = time.time()
    execution_time = end_time - start_time
    
    print(f"\nOptimal rate: {optimal_rate:.7f}")
    print(f"Final profit: {calculate_financials(optimal_rate):.2f}")
    print(f"\nInitial rate: {initial_rate}")
    print(f"Initial profit: {calculate_financials(initial_rate):.2f}")
    
    print(f"\nExecution time (excluding warm-up): {execution_time:.6f} seconds")

    # Multiple runs for average timing
    print("\nRunning 10 iterations for average timing...")
    times = []
    for _ in range(10):
        start = time.time()
        goal_seek(target_profit)
        end = time.time()
        times.append(end - start)
    
    avg_time = sum(times) / len(times)
    print(f"Average execution time over 10 runs: {avg_time:.6f} seconds")

if __name__ == "__main__":
    main()
