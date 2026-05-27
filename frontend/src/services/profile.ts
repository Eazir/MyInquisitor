import api from './api';

export interface UpdateProfileInput {
  full_name: string;
  email: string;
}

export interface ChangePasswordInput {
  current_password: string;
  new_password: string;
}

export interface UpdateProfileResponse {
  user: {
    id: string;
    email: string;
    full_name: string;
    role: string;
  };
}

export async function updateProfile(input: UpdateProfileInput): Promise<UpdateProfileResponse> {
  const { data } = await api.put('/profile', input);
  return data.data;
}

export async function changePassword(input: ChangePasswordInput): Promise<void> {
  await api.put('/profile/password', input);
}
