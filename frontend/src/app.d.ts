// See https://svelte.dev/docs/kit/types#app.d.ts

import type { Session, SupabaseClient, User } from "@supabase/supabase-js";

// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			supabase: SupabaseClient;
			safeGetSession(): Promise<{
				session: Session | null;
				user: User | null;
			}>;
			session: Session | null;
			getSession(): Promise<Session | null>;
			user: User | null;
			getUser(): Promise<User | null>;
		}
		interface PageData {
			supabase: SupabaseClient;
			session: Session | null;
			user: User | null;
		}
		// interface PageState {}
		// interface Platform {}
	}
	type Device = 'mobile' | 'desktop';
}

export {};
