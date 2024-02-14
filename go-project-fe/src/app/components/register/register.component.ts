import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrl: './register.component.scss'
})
export class RegisterComponent {
  registerData: any = {
    username: '',
    password: '',
    firstName: '',
    lastName: '',
    email: '',
    description: '',
    role: '' // You can set a default role if needed
  };
  errorMessage: string = '';

  constructor(private http: HttpClient, private router: Router) { }

  register(): void {
    this.http.post<any>(`http://127.0.0.1:3000/api/register`, this.registerData)
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