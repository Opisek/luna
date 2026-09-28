type RegistrationInvite = {
  id: string;
  author: string;
  email: string;
  created_at: Date;
  expires_at: Date;
  code: string;
}