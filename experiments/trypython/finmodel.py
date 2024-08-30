import numpy as np
from scipy.optimize import root_scalar
import pandas as pd

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
    hsi[np.where((tsn >= hsi_tsn) & (np.roll(tsn, 1) < hsi_tsn))[0]] = 1
    hsi[0] = 1 if tsn[0] >= hsi_tsn else 0
    
    overhaul = np.zeros(num_years)
    overhaul[np.where((tsn >= overhaul_tsn) & (np.roll(tsn, 1) < overhaul_tsn))[0]] = 1
    overhaul[0] = 1 if tsn[0] >= overhaul_tsn else 0
    
    total_cost = hsi * hsi_cost + overhaul * overhaul_cost
    total_profit = total_revenue - total_cost
    cumulative_profit = np.cumsum(total_profit)
    
    return cumulative_profit[-1]

def goal_seek(target_profit, initial_guess=320):
    def objective(rate):
        return calculate_financials(rate) - target_profit
    
    result = root_scalar(objective, x0=initial_guess, x1=initial_guess*1.1, method='secant', xtol=1e-8)
    return result.root

def main():
    target_profit = 3_000_000
    initial_rate = 320

    optimal_rate = goal_seek(target_profit, initial_rate)
    print(f"Optimal rate: {optimal_rate:.7f}")
    
    final_profit = calculate_financials(optimal_rate)
    print(f"Final profit: {final_profit:.2f}")

    # Verify with initial rate
    initial_profit = calculate_financials(initial_rate)
    print(f"\nInitial rate: {initial_rate}")
    print(f"Initial profit: {initial_profit:.2f}")

    # Create DataFrame for detailed view (optional)
    def create_dataframe(rate):
        num_years = 10
        years = np.arange(1, num_years + 1)
        tsn = 100 + 450 * years
        escalated_rate = rate * (1 + 5/100) ** (years - 1)
        engine_revenue = 450 * escalated_rate
        aic_revenue = engine_revenue * 10 / 100
        total_revenue = engine_revenue + aic_revenue
        
        hsi = np.zeros(num_years)
        hsi[np.where((tsn >= 1000) & (np.roll(tsn, 1) < 1000))[0]] = 1
        hsi[0] = 1 if tsn[0] >= 1000 else 0
        
        overhaul = np.zeros(num_years)
        overhaul[np.where((tsn >= 3000) & (np.roll(tsn, 1) < 3000))[0]] = 1
        overhaul[0] = 1 if tsn[0] >= 3000 else 0
        
        total_cost = hsi * 50000 + overhaul * 100000
        total_profit = total_revenue - total_cost
        cumulative_profit = np.cumsum(total_profit)
        
        return pd.DataFrame({
            'Year': years,
            'TSN': tsn,
            'Escalated Rate': escalated_rate,
            'Engine Revenue': engine_revenue,
            'AIC Revenue': aic_revenue,
            'Total Revenue': total_revenue,
            'HSI': hsi,
            'Overhaul': overhaul,
            'Total Cost': total_cost,
            'Total Profit': total_profit,
            'Cumulative Profit': cumulative_profit
        })

    # Uncomment to print detailed DataFrame
    # print("\nDetailed DataFrame for optimal rate:")
    # print(create_dataframe(optimal_rate))

if __name__ == "__main__":
    main()
