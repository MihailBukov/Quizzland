import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

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

  constructor(private http: HttpClient) { }

  login() {
    // Send login request to backend
    this.http.post<any>('http://127.0.0.1:3000/api/login', this.loginData)
      .subscribe(
        response => {
          // Handle successful login
          console.log('Login successful:', response);
          // You can perform additional actions after successful login
        },
        error => {
          // Handle login error
          console.error('Login error:', error);
          this.errorMessage = 'Invalid username or password';
        }
      );
  }
}