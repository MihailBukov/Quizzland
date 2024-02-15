import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

interface LoginRequest {
  username: string;
  password: string;
}

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {
  loginData: LoginRequest = {
    username: '',
    password: ''
  };
  errorMessage: string = '';

  constructor(private router: Router, private http: HttpClient) { }

  login() {
    this.http.post<any>('http://127.0.0.1:3000/api/login', this.loginData, {
      withCredentials: true
    })
      .subscribe(
        response => {
          this.router.navigate(['marketplace']);
        },
        error => {
          this.errorMessage = 'Invalid username or password';
        }
      );
  }
}