export interface LoginResponse {
    token?: string; // If JWT is returned
    user?: any;
    // If it sets cookie, response might be empty or success message
}

export interface LocalLoginRequest {
    email: string;
    password: string;
}
