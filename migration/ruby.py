
import pandas as pd
from datetime import date, datetime, timedelta

processStartTime = datetime.now()

def calculate_contract_periods(start_date, end_date):
    # Ensure start_date and end_date are datetime objects
    if isinstance(start_date, str):
        start_date = datetime.strptime(start_date, '%m/%d/%Y')
    if isinstance(end_date, str):
        end_date = datetime.strptime(end_date, '%m/%d/%Y')

    periods = []

    if start_date.day <= 14:
        first_period_end = datetime(start_date.year + 1, start_date.month, 1) - timedelta(days=1)
    else:
        first_period_end = datetime(start_date.year + 1, start_date.month + 1, 1) - timedelta(days=1)
        first_period_end = first_period_end.replace(day=1) + timedelta(days=31)
        first_period_end = first_period_end.replace(day=1) - timedelta(days=1)

    if first_period_end > end_date:
        first_period_end = end_date
    days = (first_period_end - start_date).days + 1
    if days >= 300:
        periods.append((start_date.strftime('%m/%d/%Y'), first_period_end.strftime('%m/%d/%Y'), days))
    start_of_next_period = first_period_end + timedelta(days=1)

    while start_of_next_period.year < end_date.year:
        next_period_end = datetime(start_of_next_period.year + 1, start_of_next_period.month, 1) - timedelta(days=1)
        days = (next_period_end - start_of_next_period).days + 1
        if days >= 300:
            periods.append((start_of_next_period.strftime('%m/%d/%Y'), next_period_end.strftime('%m/%d/%Y'), days))
        start_of_next_period = next_period_end + timedelta(days=1)

    if start_of_next_period <= end_date:
        days = (end_date - start_of_next_period).days + 1
        if days >= 300:
            periods.append((start_of_next_period.strftime('%m/%d/%Y'), end_date.strftime('%m/%d/%Y'), days))

    return periods

# Function to get the minimum warranty date
def calculate_earliest_warranty_expiration(warranty_exp_date, warranty_exp_hours, au_hours_per_day, run_out_date):
    # Ensure dates are datetime objects
    if isinstance(warranty_exp_date, str):
        warranty_exp_date = datetime.strptime(warranty_exp_date, '%m/%d/%Y')
    if isinstance(run_out_date, str):
        run_out_date = datetime.strptime(run_out_date, '%m/%d/%Y')

    # Convert Warranty Expiration Hours to Date
    days_from_hours = warranty_exp_hours / au_hours_per_day
    expiration_date_from_hours = run_out_date + timedelta(days=days_from_hours)

    # Determine the earliest expiration date
    earliest_expiration_date = min(warranty_exp_date, expiration_date_from_hours)

    return earliest_expiration_date

# Function to calculate warranty days
def calculate_warranty_days(row, first_run_rate_switch_date):
    start_date = row['RunoutStartDate']
    end_date = row['MMCContractEndDate']
    first_run_rate_switch_date = first_run_rate_switch_date

    if first_run_rate_switch_date < start_date:
        return 0
    elif first_run_rate_switch_date > end_date:
        return (end_date - start_date).days + 1
    else:
        return (first_run_rate_switch_date - start_date).days + 1

# Function to calculate first run rate days
def calculate_first_run_rate_days(row, second_run_rate_switch_date):
    start_date = row['RunoutStartDate']
    end_date = row['MMCContractEndDate']
    second_run_rate_switch_date = second_run_rate_switch_date

    if second_run_rate_switch_date < start_date:
        return 0
    elif second_run_rate_switch_date > end_date:
        return (end_date - start_date).days + 1
    else:
        return (second_run_rate_switch_date - start_date).days + 1

# Function to calculate days within a period
def calculate_days_within_period(start_date, end_date, period_start, period_end):
    if start_date > period_end or end_date < period_start:
        return 0
    actual_start = max(start_date, period_start)
    actual_end = min(end_date, period_end)
    return (actual_end - actual_start).days + 1

###########################################################
# Main Execution
# Input variables, these will come from Elevate when goalSeek is called.
scenarioId = '9076'
targetValue = 5000000
# Parameters for targetKpi - CumulativeTotalRevenue, 
targetKpi = 'CumulativeTotalRevenue'
# Parameters for goalseekUsingParameter - warranty_rate,first_run_rate,second_run_rate,third_run_rate
goalseekUsingParameter = 'first_run_rate'
###########################################################

# Process the runout calculation 
# Optional parameters are used for goal seek calculation
def process_runout(warranty_rate = 0.0, first_run_rate = 0.0, second_run_rate = 0.0, third_run_rate = 0.0):

    # Used for calculating the AUHours per day from annual AUHours
    numOfDaysInAYear = 365
    auHours = 480.0
    flightHoursMinimum = 150
    managementFees = 15.0
    aicFees = 20.0
    trustLoadFees = 2.98
    buyIn = 1352291.05
    contractStartDate = '2024-01-01T00:00:00'
    contractEndDate = '2033-12-31T23:59:59'
    
    # Identifying number of distinct assets/Engines
    rate_trend_values = [
        1, 1.0875, 1.18265625, 1.286138671875, 1.39867580566406, 1.52105993865967,
        1.65415268329239, 1.79889104308047, 1.95629400935001, 2.12746973516814,
        2.31362333699535, 2.51606537898244, 2.73622109964341, 2.97564044586221,
        3.23600898487515, 3.51915977105172, 3.82708625101875, 4.16195629798289,
        4.52612747405639, 4.92216362803633, 5.3528529454895, 5.82122757821983,
        6.33058499131407, 6.88451117805405, 7.48690590613378, 8.14201017292048,
        8.85443606305102, 9.62919921856799, 10.4717541501927, 11.3880326383345,
        12.3844854941888, 13.4681279749303, 14.6465891727367, 15.9281657253512,
        17.3218802263194, 18.8375447461224
    ]
    engineValues = [1085718,1085719]

    # Convert the strings to datetime objects
    contractStartDate = datetime.fromisoformat(contractStartDate.replace("Z", "+00:00"))
    contractEndDate = datetime.fromisoformat(contractEndDate.replace("Z", "+00:00"))

    # start calculating the mmc contract anniversary date periods
    start_datetime = datetime.combine(contractStartDate, datetime.min.time())
    end_datetime = datetime.combine(contractEndDate, datetime.min.time())
    contractDatePeriods = calculate_contract_periods(start_datetime, end_datetime)

    dfMmcContractDatesOriginal = pd.DataFrame(contractDatePeriods, columns=['MMCContractStartDate', 'MMCContractEndDate', 'NumOfDays'])
    dfMmcContractDatesOriginal['ContractYearNumber'] = range(1, len(dfMmcContractDatesOriginal)+1)
    # Convert date columns to datetime
    dfMmcContractDatesOriginal['MMCContractStartDate'] = pd.to_datetime(dfMmcContractDatesOriginal['MMCContractStartDate'], format='%m/%d/%Y')
    dfMmcContractDatesOriginal['MMCContractEndDate'] = pd.to_datetime(dfMmcContractDatesOriginal['MMCContractEndDate'], format='%m/%d/%Y')
    # Current date
    runoutDate = datetime.now()
    # Find the row where the current date falls
    currentRow = dfMmcContractDatesOriginal[(dfMmcContractDatesOriginal['MMCContractStartDate'] <= runoutDate) & 
                                    (dfMmcContractDatesOriginal['MMCContractEndDate'] >= runoutDate)].index[0]
    dfMmcContractDates = dfMmcContractDatesOriginal.iloc[currentRow:]
    
    # Create the new dataframe
    newData = {
        'RunoutStartDate': [runoutDate.strftime('%m/%d/%Y')] + 
            dfMmcContractDatesOriginal.loc[currentRow + 1:,
                                'MMCContractStartDate'].dt.strftime('%m/%d/%Y').tolist(),
        'RunoutEndDate': dfMmcContractDatesOriginal.loc[currentRow:,
                                'MMCContractEndDate'].dt.strftime('%m/%d/%Y').tolist(),
        'NumofRunoutDays': [(dfMmcContractDatesOriginal.loc[currentRow, 'MMCContractEndDate'].dayofyear - runoutDate.timetuple().tm_yday) + 1] +
            dfMmcContractDatesOriginal.loc[currentRow + 1:,
                                'NumOfDays'].tolist()
    }

    dfRunoutDates = pd.DataFrame(newData)
    dfMmcContractDates = pd.concat([dfMmcContractDates.reset_index(drop=True), dfRunoutDates], axis=1)
    numOfRowsInContractsTable = len(dfMmcContractDates)
    dfMmcContractDates['RunoutStartDate'] = pd.to_datetime(dfMmcContractDates['RunoutStartDate'], format='%m/%d/%Y')
    dfMmcContractDates['RunoutEndDate'] = pd.to_datetime(dfMmcContractDates['RunoutEndDate'], format='%m/%d/%Y')
    dfMmcContractDates['NumofRunoutDays'] = dfMmcContractDates['NumofRunoutDays']
    dfMmcContractDates['AUHours'] = auHours
    dfMmcContractDates['AUHoursPerDay'] = auHours/numOfDaysInAYear


    # Calculate Run Rate Days - Start
    # TEMPORARY CODE FOR SETTING SWITCH MILESTONE DATES
    # As of now assigning static dates, later generate them dynamically based on Events or get it from database (loop thru numofengines)
    globals()[f'engine1WarrantyDate'] = datetime.strptime('08/26/2024', '%m/%d/%Y')
    globals()[f'engine1FirstRunRateSwitchDate'] = datetime.strptime('05/30/2026', '%m/%d/%Y')
    # engine1FirstRunRateSwitchDate - it is the date where first run rate ends
    globals()[f'engine1SecondRunRateSwitchDate'] = datetime.strptime('05/30/2028', '%m/%d/%Y')
    globals()[f'engine1ThirdRunRateSwitchDate'] = datetime.strptime('09/27/2046', '%m/%d/%Y')

    globals()[f'engine2WarrantyDate'] = datetime.strptime('08/26/2024', '%m/%d/%Y')
    globals()[f'engine2FirstRunRateSwitchDate'] = datetime.strptime('05/30/2038', '%m/%d/%Y')
    globals()[f'engine2SecondRunRateSwitchDate'] = datetime.strptime('05/30/2038', '%m/%d/%Y')
    globals()[f'engine2ThirdRunRateSwitchDate'] = datetime.strptime('09/27/2046', '%m/%d/%Y')

    # Create Dynamic variables to store engine id's
    for i, value in enumerate(engineValues):
        globals()[f'engine{i+1}'] = f'Engine{i+1}FHRevenue'
        print(globals()[f'engine{i+1}'])
        # Initialize the run rate columns
        dfMmcContractDates[f'Engine{i+1}WarrantyRateDays'] = 0
        dfMmcContractDates[f'Engine{i+1}FirstRunRateDays'] = 0
        dfMmcContractDates[f'Engine{i+1}SecondRunRateDays'] = 0
        dfMmcContractDates[f'Engine{i+1}ThirdRunRateDays'] = 0

    # Loop through number of engines to calculate the switch milestone dates for each engine
    for i, value in enumerate(engineValues):
        # Calculate the days for each period
        for index, row in dfMmcContractDates.iterrows():
            period_start = row['RunoutStartDate']
            period_end = row['RunoutEndDate']    

            # WarrantyRateDays
            dfMmcContractDates.at[index, f'Engine{i+1}WarrantyRateDays'] = calculate_days_within_period(period_start, globals()[f'engine{i+1}WarrantyDate'], period_start, period_end)
            # FirstRunRateDays
            dfMmcContractDates.at[index, f'Engine{i+1}FirstRunRateDays'] = calculate_days_within_period(globals()[f'engine{i+1}WarrantyDate'] + timedelta(days=1), globals()[f'engine{i+1}FirstRunRateSwitchDate'], period_start, period_end)
            # SecondRunRateDays
            dfMmcContractDates.at[index, f'Engine{i+1}SecondRunRateDays'] = calculate_days_within_period(globals()[f'engine{i+1}FirstRunRateSwitchDate'] + timedelta(days=1), globals()[f'engine{i+1}SecondRunRateSwitchDate'], period_start, period_end)
            # ThirdRunRateDays
            if period_end <= globals()[f'engine{i+1}ThirdRunRateSwitchDate']:
                dfMmcContractDates.at[index, f'Engine{i+1}ThirdRunRateDays'] = calculate_days_within_period(globals()[f'engine{i+1}SecondRunRateSwitchDate'] + timedelta(days=1), globals()[f'engine{i+1}ThirdRunRateSwitchDate'], period_start, period_end)
            else:
                dfMmcContractDates.at[index, f'Engine{i+1}ThirdRunRateDays'] = (period_end - period_start).days + 1
        # Calculate TotalDays
        dfMmcContractDates[f'Engine{i+1}TotalDays'] = dfMmcContractDates[[f'Engine{i+1}WarrantyRateDays', f'Engine{i+1}FirstRunRateDays', f'Engine{i+1}SecondRunRateDays', f'Engine{i+1}ThirdRunRateDays']].sum(axis=1)
    # Calculate Run Rate Days - End

    dfMmcContractDates['WarrantyRate'] = warranty_rate
    if not first_run_rate or first_run_rate == 0.0:
        first_run_rate = warranty_rate
    dfMmcContractDates['FirstRunRate'] = first_run_rate
    if not second_run_rate or second_run_rate == 0.0:
        second_run_rate = first_run_rate
    dfMmcContractDates['SecondRunRate'] = second_run_rate
    if not third_run_rate or third_run_rate == 0.0:
        third_run_rate = second_run_rate
    dfMmcContractDates['ThirdRunRate'] = third_run_rate

    # Add the rateTrendColumn to Main table
    dfMmcContractDates['RateTrend'] = rate_trend_values[:numOfRowsInContractsTable]

    # Calculate Engine Flight Hours Revenue (FHRevenue) - loop through number of engines
    for i, value in enumerate(engineValues):
        dfMmcContractDates[f'Engine{i+1}WarrantyCalc'] = (dfMmcContractDates[f'Engine{i+1}WarrantyRateDays'] * dfMmcContractDates['WarrantyRate'])
        dfMmcContractDates[f'Engine{i+1}FirstRunRateCalc'] = (dfMmcContractDates[f'Engine{i+1}FirstRunRateDays'] * dfMmcContractDates['FirstRunRate'])
        dfMmcContractDates[f'Engine{i+1}SecondRunRateCalc'] = (dfMmcContractDates[f'Engine{i+1}SecondRunRateDays'] * dfMmcContractDates['SecondRunRate'])
        dfMmcContractDates[f'Engine{i+1}ThirdRunRateCalc'] = (dfMmcContractDates[f'Engine{i+1}ThirdRunRateDays'] * dfMmcContractDates['ThirdRunRate'])
        dfMmcContractDates[f'Engine{i+1}Rates'] = (
            dfMmcContractDates[f'Engine{i+1}WarrantyCalc'] + dfMmcContractDates[f'Engine{i+1}FirstRunRateCalc'] +
            dfMmcContractDates[f'Engine{i+1}SecondRunRateCalc'] + dfMmcContractDates[f'Engine{i+1}ThirdRunRateCalc']
        )
        dfMmcContractDates[f'Engine{i+1}EscalatedRate'] = dfMmcContractDates[f'Engine{i+1}Rates'] * dfMmcContractDates['RateTrend']
        dfMmcContractDates[f'Engine{i+1}FHUtilization'] = dfMmcContractDates['AUHoursPerDay'] * dfMmcContractDates[f'Engine{i+1}TotalDays']
        dfMmcContractDates[f'Engine{i+1}Shortfall'] = dfMmcContractDates.apply(
            lambda row: flightHoursMinimum - row[f'Engine{i+1}FHUtilization'] if row[f'Engine{i+1}FHUtilization'] < flightHoursMinimum else 0,
            axis=1
        )
        dfMmcContractDates[f'Engine{i+1}FHRevenue'] = (
            dfMmcContractDates[f'Engine{i+1}EscalatedRate'] * dfMmcContractDates['AUHoursPerDay']
        )

    # Identify columns dynamically based on a pattern
    patternFHRevenue = 'FHRevenue'
    patternFHRevenueColumns = [col for col in dfMmcContractDates.columns if patternFHRevenue in col]

    # Add the identified columns and store the result in a new column
    dfMmcContractDates['TotalFHRevenue'] = dfMmcContractDates[patternFHRevenueColumns].sum(axis=1)
    dfMmcContractDates['MgmtFeeRevenue'] = (dfMmcContractDates['TotalFHRevenue'] * (managementFees/100))
    dfMmcContractDates['AICRevenue'] = dfMmcContractDates['TotalFHRevenue'] * (1-(managementFees/100)) * (aicFees/100)
    dfMmcContractDates['TrustLoadRevenue'] = dfMmcContractDates['TotalFHRevenue'] * (1-(managementFees/100)) * (trustLoadFees/100)
    dfMmcContractDates['BuyIn'] = 0.0
    dfMmcContractDates.at[0, 'BuyIn'] = buyIn
    dfMmcContractDates['TrustRevenue'] = ( dfMmcContractDates['TotalFHRevenue'] -
        (dfMmcContractDates['MgmtFeeRevenue'] + dfMmcContractDates['AICRevenue'] + dfMmcContractDates['TrustLoadRevenue'] + dfMmcContractDates['BuyIn'])
    )
    dfMmcContractDates['TotalRevenue'] = (
        dfMmcContractDates['MgmtFeeRevenue'] + dfMmcContractDates['AICRevenue'] + dfMmcContractDates['TrustLoadRevenue'] +
        dfMmcContractDates['BuyIn'] + dfMmcContractDates['TrustRevenue']
    )

    dfMmcContractDates['CumulativeTotalRevenue'] = dfMmcContractDates['TotalRevenue'].cumsum()

    #pd.set_option('display.max_columns', None) #display all columns
    print(dfMmcContractDates)
    dfMmcContractDates.to_csv('runout_analysis_ruby_results_v2.csv', index=False)
    print("\nResults saved to 'runout_analysis_ruby_results_v2.csv'")
    return dfMmcContractDates


warranty_rate = 243.6
first_run_rate = 255.13
second_run_rate = 255.13
third_run_rate = 255.13
process_runout(warranty_rate,first_run_rate,second_run_rate,third_run_rate)

# Process execution end time
processEndTime = datetime.now()
# Calculating overall process execution time
executionTime = processEndTime - processStartTime
print(f"Overall Execution Time: {executionTime}")