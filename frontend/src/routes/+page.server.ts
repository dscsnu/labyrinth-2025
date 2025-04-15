import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

// @ts-ignore
export const load: PageServerLoad = async ({ locals }) => {
    if (locals.user) {
        redirect(303, '/new')
    }

    return;
}