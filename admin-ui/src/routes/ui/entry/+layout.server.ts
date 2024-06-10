import type { LayoutServerLoad } from '../$types';

export const load: LayoutServerLoad = async ({ fetch }) => {
      const res = await fetch('http://localhost:6701/api/v1/user_schema',{
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
      });
        const schema = await res.json();
        return JSON.parse(schema);
}