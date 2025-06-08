export interface StopSummary {
  stop_id: number;
  stop_name: string;
  lines: string[];
}

export interface StopDetails {
  stop_id: number;
  stop_name: string;
  image_url: string;
  lines: string[];
  departures: StopDeparture[];
}

export interface StopDeparture {
  tram_id: string;
  line: string;
  destination: string;
  arrival_time: string;
}
