using System;
using System.Diagnostics;  // For Stopwatch
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;  // For JsonPropertyName
using System.Threading.Tasks;

public class GoalSeekRequest
{
    public int NumYears { get; set; }
    public int AuHours { get; set; }
    public int InitialTSN { get; set; }
    public double RateEscalation { get; set; }
    public int Aic { get; set; }
    public int HsiTsn { get; set; }
    public int OverhaulTSN { get; set; }
    public int HsiCost { get; set; }
    public int OverhaulCost { get; set; }
    public int TargetProfit { get; set; }
    public double InitialRate { get; set; }
}

public class GoalSeekResponse
{
    [JsonPropertyName("finalCumulativeProfit")]
    public double FinalCumulativeProfit { get; set; }

    [JsonPropertyName("iterations")]
    public int Iterations { get; set; }

    [JsonPropertyName("optimalWarrantyRate")]
    public double OptimalWarrantyRate { get; set; }
}

public class Program
{
    private static readonly HttpClient client = new HttpClient();
    private const string ApiUrl = "http://localhost:8080/goalseek";

    public static async Task Main()
    {
        var request = new GoalSeekRequest
        {
            NumYears = 10,
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

        try
        {
            var stopwatch = Stopwatch.StartNew(); // Start the stopwatch
            var result = await CallGoalSeekApiAsync(request);
            stopwatch.Stop(); // Stop the stopwatch

            Console.WriteLine($"Response Time: {stopwatch.ElapsedMilliseconds} ms");

            if (result != null)
            {
                Console.WriteLine($"Optimal Warranty Rate: {result.OptimalWarrantyRate}");
                Console.WriteLine($"Iterations: {result.Iterations}");
                Console.WriteLine($"Final Cumulative Profit: {result.FinalCumulativeProfit}");
            }
            else
            {
                Console.WriteLine("API response was null or deserialization failed.");
            }
        }
        catch (HttpRequestException e)
        {
            Console.WriteLine($"Error calling API: {e.Message}");
        }
    }

    public static async Task<GoalSeekResponse?> CallGoalSeekApiAsync(GoalSeekRequest request)
    {
        var json = JsonSerializer.Serialize(request);
        var content = new StringContent(json, Encoding.UTF8, "application/json");

        var response = await client.PostAsync(ApiUrl, content);
        response.EnsureSuccessStatusCode();

        var responseBody = await response.Content.ReadAsStringAsync();
        Console.WriteLine("Response from API: " + responseBody);

        return JsonSerializer.Deserialize<GoalSeekResponse>(responseBody);
    }
}
