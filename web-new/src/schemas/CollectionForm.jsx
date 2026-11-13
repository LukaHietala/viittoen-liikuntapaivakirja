import { createForm, Field, Form, reset, setInput, getInput } from "@formisch/solid";
import { createEffect } from "solid-js";
import * as v from "valibot";

const collectionSchema = v.pipe(
		v.object({
				title: v.pipe(
						v.string(),
						v.minLength(3, "Nimessä tulee olla vähintään 3 merkkiä"),
						v.maxLength(50, "Nimessä voi olla enintään 50 merkkiä"),
				),
				description: v.pipe(
						v.string(),
						v.minLength(3, "Kuvauksen tulee olla vähintään 3 merkkiä"),
						v.maxLength(200, "Kuvauksessa voi olla enintään 200 merkkiä"),
				),
				goalPoints: v.pipe(
						v.toNumber(),
						v.number(),
						v.integer(),
						v.toMinValue(1)
				),
				startDate: v.pipe(v.string(), v.isoDate()),
				endDate: v.pipe(v.string(), v.isoDate()),
		}),
		v.forward(
				v.partialCheck(
						[["startDate"], ["endDate"]],
						(input) => {
								return input.endDate > input.startDate;
						},
						"Aloituspäivän tulee olla ennen lopetuspäivää",
				),
				["endDate"],
		),
);

function CollectionForm(props) {
		const collection = props.collection;
		const cancel = () => props.cancel;

		const handleSubmit = async (data) => {
				await props.onSubmit(data);
		};

		const collectionForm = createForm({
				schema: collectionSchema,
		});

		createEffect(() => {
				if (collection) {
						reset(collectionForm, {
								initialInput: collection(),
						});
				}
		});

		return (
				<Form
						of={collectionForm}
						onSubmit={handleSubmit}
						class="form"
				>				
						<Field
								of={collectionForm}
								path={["title"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Nimi:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="text"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>

						<Field
								of={collectionForm}
								path={["description"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Kuvaus:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="textarea"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						
						<Field
								of={collectionForm}
								path={["startDate"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Aloituspäivä:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="date"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						<Field
								of={collectionForm}
								path={["endDate"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Lopetuspäivä:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="date"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						
						<Field
								of={collectionForm}
								path={["goalPoints"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Tavoite pisteet:</label>
												<br/>
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="number"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						
						<button
								class="btn-secondary"
								type="button"
								onClick={() => cancel()}
						>
								Peruuta
						</button>
						<button
								class="btn-primary"
								type="submit"
						>
								Tallenna
						</button>
				</Form>
		);
}

export default CollectionForm;