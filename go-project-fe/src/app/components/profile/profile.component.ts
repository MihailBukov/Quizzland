import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

interface ProfileDto {
  id: number;
  username: string;
  firstName: string;
  lastName: string;
  email: string;
  description: string;
  balance: number;
  quizzes: QuizDto[];
}

interface QuizDto {
  id: number;
  name: string;
  description: string;
  owner: string;
}

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrl: './profile.component.scss'
})
export class ProfileComponent implements OnInit {

  profile: ProfileDto;

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    this.http.get<ProfileDto>('/api/profile').subscribe(data => {
      this.profile = data;
    });
  }

  updateProfile() {
    this.http.post('/api/profile/update', this.profile).subscribe(() => {
      console.log('Profile updated successfully');
    });
  }

}