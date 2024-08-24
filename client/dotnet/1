using System;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using System.Text.Json.Serialization;
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
            var result = await CallGoalSeekApiAsync(request);
            Console.WriteLine($"Optimal Warranty Rate: {result.OptimalWarrantyRate}");
            Console.WriteLine($"Iterations: {result.Iterations}");
            Console.WriteLine($"Final Cumulative Profit: {result.FinalCumulativeProfit}");
        }
        catch (HttpRequestException e)
        {
            Console.WriteLine($"Error calling API: {e.Message}");
        }
        catch (Exception e)
        {
            Console.WriteLine($"An unexpected error occurred: {e.Message}");
        }
    }


    private static async Task<GoalSeekResponse?> CallGoalSeekApiAsync(GoalSeekRequest request)
{
    // Serialize the request object to JSON format
    var json = JsonSerializer.Serialize(request);

    // Create the HTTP content using the JSON string
    var content = new StringContent(json, Encoding.UTF8, "application/json");

    // Send the POST request to the API
    var response = await client.PostAsync(ApiUrl, content);

    // Ensure that the request was successful
    response.EnsureSuccessStatusCode();

    // Read the response body as a string
    var responseBody = await response.Content.ReadAsStringAsync();

    // Log the raw JSON response from the API
    Console.WriteLine("Response from API: " + responseBody);

    // Deserialize the JSON response into the GoalSeekResponse object
    var result = JsonSerializer.Deserialize<GoalSeekResponse>(responseBody);

    // Return the deserialized object or null if deserialization failed
    return result;
}

}
