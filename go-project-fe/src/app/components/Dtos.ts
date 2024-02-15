export interface ProfileDto {
    id: number;
    username: string;
    firstName: string;
    lastName: string;
    email: string;
    description: string;
    balance: number;
    quizzes: QuizDto[];
  }
  
  export interface QuizDto {
    id: number;
    name: string;
    description: string;
    owner: string;
  }