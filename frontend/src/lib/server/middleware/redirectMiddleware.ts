import { redirect, type Handle } from '@sveltejs/kit';
import { TOKEN_NAME as TEAM_TOKEN_NAME, type ITeamData } from '$lib/stores/TeamStore';

export const authGuardMiddleware: Handle = async ({ event, resolve }) => {
    const currentPath = event.url.pathname;

    if (currentPath.startsWith('/api')) {
        return resolve(event);
    }

    // Ensure user is authenticated
    const user = event.locals.user;
    if (!user) {
        // Allow access only to landing/login page
        if (currentPath !== '/') {
            return redirect(303, '/');
        }
        return resolve(event);
    }

    // Ensure user has a team
    const teamCookie = event.cookies.get(TEAM_TOKEN_NAME);
    if (!teamCookie) {
        return redirect(303, '/newplayer')
    }

    // Extract team information
    let teamData: ITeamData;
    try {
        teamData = JSON.parse(teamCookie);
    } catch (e) {
        console.error(`Failed parsing team cookie: ${e}`);
    }

    // Redirect to /finish if gameFinished
    // if (game has finished %% currentPath !== '/finish') {
    //     return redirect(303, '/finish');
    // }


    // User authenticated but on landing page
    if (currentPath === '/') {
        // Redirect to team page if they have a team, otherwise to newplayer
        if (teamData!.id) {
            return redirect(303, '/team');
        } else {
            return redirect(303, '/newplayer');
        }
    }

    // User authenticated but no team
    if (!teamData!.id) {
        // If not on newplayer page, redirect there
        if (currentPath !== '/newplayer') {
            return redirect(303, '/newplayer');
        }
        return resolve(event);
    }

    // User has team but tries to access newplayer
    if (teamData!.id && currentPath === '/newplayer') {
        return redirect(303, '/team');
    }


    // Game flow controls
    if (teamData!.id) {

        // if (Game has started && team.isReady) {
        //      return redirect(303, '/game');
        // }
        if (currentPath === '/game') {
            // Only allow access if team is ready and game has started
            if (!teamData!.allReady) {
                // Redirect to team page if not ready
                return redirect(303, '/team');
            }

            // if (Game has not started yet){
            //     return redirect(303, '/team');
            // }
        }

        if (currentPath === '/finish') {
            // Game ended check comes here
        }
    }

    // All checks passed, allow access to the requested page
    return resolve(event);
}