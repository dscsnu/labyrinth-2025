import { error, json } from '@sveltejs/kit';
import jwt from 'jsonwebtoken';
import { JWT_SESSION_KEY } from "$env/static/private";

export async function POST({ request }) {
    console.log('Received request to decrypt token');
    const data = await request.json();

    try {
        console.log('Received data:', data.token);
        if (!data.token) {
            throw error(400, 'Token is required');
        }

        const decoded = jwt.verify(data.token, JWT_SESSION_KEY) as {
            payload: string;
            timestamp: number;
        };
        let valid = false;
        let validTime = 5; //in seconds
        if (decoded && decoded.payload && decoded.timestamp) {
            const currentTime = Math.floor(Date.now() / 1000);
            if (currentTime < decoded.timestamp + validTime) {
                valid = true;
            }
        }
        let route = "/";
        if (decoded.payload.startsWith('location:')) {
            route = "/api/scan/location";
        } 
        else if (decoded.payload.startsWith('qte:')) {
            route = "/api/scan/qte";
        }        

        return json({
            valid: valid,
            route: route,
        });
    } catch (err) {
        console.error('Error decrypting token:', err);
        return error(500, 'Failed to decrypt token');
    }
}