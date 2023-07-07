import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../environments/environment';
import { Observable, combineLatest, lastValueFrom } from 'rxjs';
import { Station, NetworkStats } from './station';

@Injectable({
  providedIn: 'root'
})
export class StationService {
  private stationsUrl = environment.baseUrl + '/app/GetStations/';
  private statsUrl = environment.baseUrl + '/app/netStats/';

  constructor(private http: HttpClient) {}

  async getStats(): Promise<NetworkStats> {
    let ns = await lastValueFrom(this.http.get<NetworkStats>(this.statsUrl));
    let obj = new NetworkStats(ns);
    console.log(ns);
    return obj;
  }

  async getStations(): Promise<Station[]> {
    let stationData = await lastValueFrom(this.http.get<Station[]>(this.stationsUrl));
    return stationData.map(item => new Station(item));
  }

  async getStation(name: string): Promise<Station | undefined> {
    let stations = await this.getStations();
    for (let station of stations) {
      if (station.name === name) {
        return station;
      }
    }
    return undefined;
  }
}
