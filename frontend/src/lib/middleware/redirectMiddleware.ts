import { redirect, type Handle } from '@sveltejs/kit';
import { TEAM_TOKEN_NAME } from '$lib/stores/TeamStore';

export const redirectMiddleware: Handle = async ({ event, resolve }) => {
    const currentPath = event.url.pathname;

    if (currentPath.startsWith('/api')) {
        return resolve(event);
    }

    const user = event.locals.user;
    
    const teamCookie = event.cookies.get(TEAM_TOKEN_NAME);
    let teamData = null;
    
    if (teamCookie) {
        try {
            teamData = JSON.parse(decodeURIComponent(teamCookie));
        } catch (e) {
            console.error('Failed to parse team cookie:', e);
        }
    }

    // User not authenticated
    if (!user) {
        // Allow access only to landing/login page
        if (currentPath !== '/') {
            return redirect(303, '/');
        }
        return resolve(event);
    }
    
    // If user has logged in and game has finished
    // if (game has finished %% currentPath !== '/finish') {
    //     return redirect(303, '/finish');
    // }

    
    // User authenticated but on landing page
    if (currentPath === '/') {
        // Redirect to team page if they have a team, otherwise to newplayer
        if (teamData?.id) {
            return redirect(303, '/team');
        } else {
            return redirect(303, '/newplayer');
        }
    }
    
    // User authenticated but no team
    if (!teamData?.id) {
        // If not on newplayer page, redirect there
        if (currentPath !== '/newplayer') {
            return redirect(303, '/newplayer');
        }
        return resolve(event);
    }

    // User has team but tries to access newplayer
    if (teamData?.id && currentPath === '/newplayer') {
        return redirect(303, '/team');
    }


    // Game flow controls
    if (teamData?.id) {

        // if (Game has started){
        //      return redirect(303, '/game');
        // }
        if (currentPath === '/game') {
            // Only allow access if team is ready and game has started
            if (!teamData.isReady) {
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

export default redirectMiddleware;