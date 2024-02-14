import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-topbar',
  templateUrl: './topbar.component.html',
  styleUrl: './topbar.component.scss'
})
export class TopbarComponent {

  constructor(private router: Router) { }

  isLoggedIn(): boolean {
    // Implement logic to check if the user is logged in
    // You can retrieve this information from cookies
    return false; // For demonstration purposes, always return true
  }

  goToLogin(): void {
    this.router.navigate(['login']);
  }

  goToRegister(): void {
    this.router.navigate(['register']);
  }

  goToMarketplace(): void {
    // Implement navigation to the marketplace component
  }

  createLobby(): void {
    // Implement logic to create a lobby
  }

  joinLobby(): void {
    // Implement logic to join a lobby
  }

  createQuiz(): void {
    // Implement logic to create a quiz
  }

  goToProfile(): void {
    // Implement navigation to the profile component
  }
}