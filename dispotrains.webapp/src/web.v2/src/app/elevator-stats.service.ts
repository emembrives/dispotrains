import { Injectable } from '@angular/core';
import { environment } from '../environments/environment';
import { HttpClient } from '@angular/common/http';
import { ElevatorStatistics }  from './elevator-stats';
import { Observable, lastValueFrom } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ElevatorStatsService {
  private elevatorUrl = environment.baseUrl + '/app/Elevator/';

  constructor(private _http: HttpClient) {}
  
  async getElevatorStats(name: string): Promise<ElevatorStatistics> {
    let resp = await lastValueFrom(this._http.get<ElevatorStatistics>(this.elevatorUrl + name));
    return new ElevatorStatistics(resp);
  }
}
