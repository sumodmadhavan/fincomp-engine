from datetime import datetime

import numpy as np
import pandas as pd


def process_runout(
    rate_trend_values,
    engine_values,
    warranty_rate=0.0,
    first_run_rate=0.0,
    second_run_rate=0.0,
    third_run_rate=0.0,
    **kwargs,
):
    # Default parameters
    params = {
        "num_of_days_in_year": 365,
        "au_hours": 480.0,
        "flight_hours_minimum": 150,
        "management_fees": 15.0,
        "aic_fees": 20.0,
        "trust_load_fees": 2.98,
        "buy_in": 1352291.05,
        "contract_start_date": "2024-01-01T00:00:00",
        "contract_end_date": "2033-12-31T23:59:59",
    }
    params.update(kwargs)

    # Convert dates to datetime
    contract_start_date = datetime.fromisoformat(
        params["contract_start_date"].replace("Z", "+00:00")
    )
    contract_end_date = datetime.fromisoformat(
        params["contract_end_date"].replace("Z", "+00:00")
    )

    # Calculate contract periods
    contract_periods = calculate_contract_periods(
        contract_start_date, contract_end_date
    )
    df_mmc_contract_dates = pd.DataFrame(
        contract_periods,
        columns=["MMCContractStartDate", "MMCContractEndDate", "NumOfDays"],
    )
    df_mmc_contract_dates["ContractYearNumber"] = range(
        1, len(df_mmc_contract_dates) + 1
    )

    # Convert date columns to datetime
    df_mmc_contract_dates["MMCContractStartDate"] = pd.to_datetime(
        df_mmc_contract_dates["MMCContractStartDate"]
    )
    df_mmc_contract_dates["MMCContractEndDate"] = pd.to_datetime(
        df_mmc_contract_dates["MMCContractEndDate"]
    )

    # Add rate trend and engine-specific columns
    df_mmc_contract_dates["RateTrend"] = rate_trend_values[: len(df_mmc_contract_dates)]
    df_mmc_contract_dates["AUHours"] = params["au_hours"]
    df_mmc_contract_dates["AUHoursPerDay"] = (
        params["au_hours"] / params["num_of_days_in_year"]
    )

    # Initialize engine-specific columns
    for i, engine_id in enumerate(engine_values):
        df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"] = 0
        df_mmc_contract_dates[f"Engine{i+1}FirstRunRateDays"] = 0
        df_mmc_contract_dates[f"Engine{i+1}SecondRunRateDays"] = 0
        df_mmc_contract_dates[f"Engine{i+1}ThirdRunRateDays"] = 0

    # Calculate run rate days (vectorized)
    for i, engine_id in enumerate(engine_values):
        warranty_date = datetime.strptime(f"08/26/2024", "%m/%d/%Y")
        first_run_rate_switch_date = datetime.strptime(f"05/30/2026", "%m/%d/%Y")
        second_run_rate_switch_date = datetime.strptime(f"05/30/2028", "%m/%d/%Y")
        third_run_rate_switch_date = datetime.strptime(f"09/27/2046", "%m/%d/%Y")

        df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"] = np.minimum(
            np.maximum(
                0,
                (warranty_date - df_mmc_contract_dates["MMCContractStartDate"]).dt.days,
            ),
            df_mmc_contract_dates["NumOfDays"],
        )
        df_mmc_contract_dates[f"Engine{i+1}FirstRunRateDays"] = np.minimum(
            np.maximum(
                0,
                (
                    first_run_rate_switch_date
                    - df_mmc_contract_dates["MMCContractStartDate"]
                ).dt.days,
            )
            - df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"],
            df_mmc_contract_dates["NumOfDays"]
            - df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"],
        )
        df_mmc_contract_dates[f"Engine{i+1}SecondRunRateDays"] = np.minimum(
            np.maximum(
                0,
                (
                    second_run_rate_switch_date
                    - df_mmc_contract_dates["MMCContractStartDate"]
                ).dt.days,
            )
            - (
                df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"]
                + df_mmc_contract_dates[f"Engine{i+1}FirstRunRateDays"]
            ),
            df_mmc_contract_dates["NumOfDays"]
            - (
                df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"]
                + df_mmc_contract_dates[f"Engine{i+1}FirstRunRateDays"]
            ),
        )
        df_mmc_contract_dates[f"Engine{i+1}ThirdRunRateDays"] = df_mmc_contract_dates[
            "NumOfDays"
        ] - (
            df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"]
            + df_mmc_contract_dates[f"Engine{i+1}FirstRunRateDays"]
            + df_mmc_contract_dates[f"Engine{i+1}SecondRunRateDays"]
        )

    # Calculate revenue (vectorized)
    df_mmc_contract_dates["WarrantyRate"] = warranty_rate
    df_mmc_contract_dates["FirstRunRate"] = first_run_rate or warranty_rate
    df_mmc_contract_dates["SecondRunRate"] = (
        second_run_rate or df_mmc_contract_dates["FirstRunRate"]
    )
    df_mmc_contract_dates["ThirdRunRate"] = (
        third_run_rate or df_mmc_contract_dates["SecondRunRate"]
    )

    for i, engine_id in enumerate(engine_values):
        df_mmc_contract_dates[f"Engine{i+1}Rates"] = (
            df_mmc_contract_dates[f"Engine{i+1}WarrantyRateDays"]
            * df_mmc_contract_dates["WarrantyRate"]
            + df_mmc_contract_dates[f"Engine{i+1}FirstRunRateDays"]
            * df_mmc_contract_dates["FirstRunRate"]
            + df_mmc_contract_dates[f"Engine{i+1}SecondRunRateDays"]
            * df_mmc_contract_dates["SecondRunRate"]
            + df_mmc_contract_dates[f"Engine{i+1}ThirdRunRateDays"]
            * df_mmc_contract_dates["ThirdRunRate"]
        )
        df_mmc_contract_dates[f"Engine{i+1}EscalatedRate"] = (
            df_mmc_contract_dates[f"Engine{i+1}Rates"]
            * df_mmc_contract_dates["RateTrend"]
        )
        df_mmc_contract_dates[f"Engine{i+1}FHUtilization"] = (
            df_mmc_contract_dates["AUHoursPerDay"] * df_mmc_contract_dates["NumOfDays"]
        )
        df_mmc_contract_dates[f"Engine{i+1}Shortfall"] = np.maximum(
            0,
            params["flight_hours_minimum"]
            - df_mmc_contract_dates[f"Engine{i+1}FHUtilization"],
        )
        df_mmc_contract_dates[f"Engine{i+1}FHRevenue"] = (
            df_mmc_contract_dates[f"Engine{i+1}EscalatedRate"]
            * df_mmc_contract_dates["AUHoursPerDay"]
        )

    # Calculate total revenue
    fh_revenue_columns = [
        col for col in df_mmc_contract_dates.columns if "FHRevenue" in col
    ]
    df_mmc_contract_dates["TotalFHRevenue"] = df_mmc_contract_dates[
        fh_revenue_columns
    ].sum(axis=1)
    df_mmc_contract_dates["MgmtFeeRevenue"] = df_mmc_contract_dates[
        "TotalFHRevenue"
    ] * (params["management_fees"] / 100)
    df_mmc_contract_dates["AICRevenue"] = (
        df_mmc_contract_dates["TotalFHRevenue"]
        * (1 - params["management_fees"] / 100)
        * (params["aic_fees"] / 100)
    )
    df_mmc_contract_dates["TrustLoadRevenue"] = (
        df_mmc_contract_dates["TotalFHRevenue"]
        * (1 - params["management_fees"] / 100)
        * (params["trust_load_fees"] / 100)
    )
    df_mmc_contract_dates["BuyIn"] = 0.0
    df_mmc_contract_dates.at[0, "BuyIn"] = params["buy_in"]
    df_mmc_contract_dates["TrustRevenue"] = df_mmc_contract_dates["TotalFHRevenue"] - (
        df_mmc_contract_dates["MgmtFeeRevenue"]
        + df_mmc_contract_dates["AICRevenue"]
        + df_mmc_contract_dates["TrustLoadRevenue"]
        + df_mmc_contract_dates["BuyIn"]
    )
    df_mmc_contract_dates["TotalRevenue"] = (
        df_mmc_contract_dates["MgmtFeeRevenue"]
        + df_mmc_contract_dates["AICRevenue"]
        + df_mmc_contract_dates["TrustLoadRevenue"]
        + df_mmc_contract_dates["BuyIn"]
        + df_mmc_contract_dates["TrustRevenue"]
    )
    df_mmc_contract_dates["CumulativeTotalRevenue"] = df_mmc_contract_dates[
        "TotalRevenue"
    ].cumsum()

    return df_mmc_contract_dates


rate_trend_values = [
    1,
    1.0875,
    1.18265625,
    1.286138671875,
    1.39867580566406,
    1.52105993865967,
    1.65415268329239,
    1.79889104308047,
    1.95629400935001,
    2.12746973516814,
]  # ... add all values
engine_values = [1085718, 1085719]
warranty_rate = 243.6
first_run_rate = 255.13
second_run_rate = 255.13
third_run_rate = 255.13

result = process_runout(
    rate_trend_values,
    engine_values,
    warranty_rate,
    first_run_rate,
    second_run_rate,
    third_run_rate,
)
print(result)
