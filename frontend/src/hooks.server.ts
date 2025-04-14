import type { Handle } from "@sveltejs/kit";
import { sequence } from "@sveltejs/kit/hooks";

import { authGuardMiddleware } from "$lib/server/middleware/redirectMiddleware";
import { handleDevice } from "$lib/server/middleware/deviceMiddleware";
import { createSupabase } from "$lib/server/middleware/supabaseMiddleware";

export const handle: Handle = sequence(createSupabase, authGuardMiddleware, handleDevice)