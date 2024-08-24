using Xunit;

public class GoalSeekApiTests
{
    [Theory]
    [InlineData(10, 2999999.9999999986, 3, 505.93820432563325)]
    [InlineData(35, 3000000.000000007, 3, 70.45631874177171)]
    [InlineData(50, 3000000.0000000005, 4, 30.397407636504852)]
    public async Task TestGoalSeekApi_WithDifferentYears(int numYears, double expectedProfit, int expectedIterations, double expectedWarrantyRate)
    {
        // Arrange
        var request = new GoalSeekRequest
        {
            NumYears = numYears,
            AuHours = 450,
            InitialTSN = 100,
            RateEscalation = 5,
            Aic = 10,
            HsiTsn = 1000,
            OverhaulTSN = 3000,
            HsiCost = 50000,
            OverhaulCost = 100000,
            TargetProfit = 3000000,
            InitialRate = 320
        };

        // Act
        var result = await Program.CallGoalSeekApiAsync(request);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(expectedProfit, result.FinalCumulativeProfit);
        Assert.Equal(expectedIterations, result.Iterations);
        Assert.Equal(expectedWarrantyRate, result.OptimalWarrantyRate);
    }
}
