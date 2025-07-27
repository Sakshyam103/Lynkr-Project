/**
 * Brand Authentication Service
 * Handles brand login, token management, and authorization
 */

import axios from 'axios';

interface BrandCredentials {
  email: string;
  password: string;
}

interface BrandUser {
  id: string;
  name: string;
  email: string;
  brandId: string;
  role: string;
}

export class AuthService {
  private static TOKEN_KEY = 'brand_auth_token';
  private static USER_KEY = 'brand_user';

  static async login(credentials: BrandCredentials): Promise<BrandUser> {
    const response = await axios.post('/api/v1/brands/login', credentials);
    const { token, user } = response.data;
    
    localStorage.setItem(this.TOKEN_KEY, token);
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    
    return user;
  }

  static logout(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
  }

  static getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  static getCurrentUser(): BrandUser | null {
    const userStr = localStorage.getItem(this.USER_KEY);
    return userStr ? JSON.parse(userStr) : null;
  }

  static isAuthenticated(): boolean {
    return !!this.getToken();
  }

  static setupAxiosInterceptors(): void {
    axios.interceptors.request.use((config) => {
      const token = this.getToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    axios.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          this.logout();
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }
}