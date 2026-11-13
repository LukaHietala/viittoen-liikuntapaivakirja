import { handleResponse } from ".";

const baseUrl = "/api/performances";

export const performanceService = {
		listByChallengeId: (id) => fetch(`${baseUrl}/challenge/${id}`).then(handleResponse),

		listOwn: () => fetch(`${baseUrl}/own`).then(handleResponse),

		create: (data) =>
		fetch(`${baseUrl}/`, {
				method: "POST",
				headers: {
						"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
		}).then(handleResponse),

		deleteByID: (id) =>
		fetch(`${baseUrl}/${id}`, {
				method: "DELETE",
		}).then(handleResponse),
};