export async function handleResponse(response) {
		if (response.ok) {
				if (response.status === 204) {
						return null;
				}
				return response.json();
		}

		const err = await response
					.json()
					.then((data) => data.error)
		// Should not happen
					.catch(() => `Error: ${response.status}`);

		throw new Error(err);
}

export { challengeService } from "./challengeService.js";
export { collectionService } from "./collectionService.js";
export { userService } from "./userService.js";
export { performanceService } from "./performanceService.js";
export { authService } from "./authService.js";