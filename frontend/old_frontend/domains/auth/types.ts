/**
 * Auth Domain Types
 */

export type AuthStatus = 'authenticated' | 'unauthenticated' | 'loading'

export interface LoginRequest {
  username: string
  password: string
  remember?: boolean
}

export interface LoginResponse {
  access_token: string
  refresh_token?: string
  expires_in: number
  token_type: 'Bearer'
}

export interface LogoutRequest {
  everywhere?: boolean
}

export interface RefreshTokenResponse {
  access_token: string
  expires_in: number
}

export type OAuthProvider = 'google' | 'facebook' | 'apple' | 'gerege'

export interface OAuthLoginRequest {
  provider: OAuthProvider
  code: string
  redirect_uri: string
}

export interface PasswordResetRequest {
  email: string
}

export interface PasswordResetConfirmRequest {
  token: string
  new_password: string
  confirm_password: string
}

export interface ChangeOrganizationRequest {
  id: number
}

export interface ChangeSystemRequest {
  system_id: number
}
