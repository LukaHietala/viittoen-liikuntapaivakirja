import { handleResponse } from ".";

const baseUrl = "/api/collections";

export const collectionService = {
		listAll: () => fetch(`${baseUrl}/`).then(handleResponse),

		listActive: () => fetch(`${baseUrl}/active`).then(handleResponse),

		findByID: (id) => fetch(`${baseUrl}/${id}`).then(handleResponse),

		create: (data) =>
		fetch(`${baseUrl}/`, {
				method: "POST",
				headers: {
						"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
		}).then(handleResponse),

		updateByID: (id, data) =>
		fetch(`${baseUrl}/${id}`, {
				method: "PUT",
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