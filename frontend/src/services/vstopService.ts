export async function getVStopData() {
  return {
    stopName: "Rynek",
    lines: [3, 18, 21, 22],
    imageUrl: "/tram.jpg", // placeholder
    balance: 26019,
    departures: [
      {
        line: "03",
        destination: "Księże Małe",
        arrivalTime: "0:16",
        odds: [],
      },
      {
        line: "18",
        destination: "Gaj",
        arrivalTime: "0:16",
        odds: [
          { label: "0:00 - 0:30", value: 1.04 },
          { label: "0:31 - 1:00", value: 1.78 },
          { label: "1:01 - 2:00", value: 2.41 },
          { label: "2:01 - 🧍‍♀️", value: 4.93 },
        ],
      },
      {
        line: "21",
        destination: "Gaj",
        arrivalTime: "0:16",
        odds: [],
      },
      {
        line: "18",
        destination: "Gaj",
        arrivalTime: "0:16",
        odds: [],
      },
    ],
  };
}
