import { createForm, Field, Form, reset, setInput, getInput } from "@formisch/solid";
import { createEffect } from "solid-js";
import * as v from "valibot";

const userSchema = v.pipe(
		v.object({
				name: v.pipe(
						v.string(),
						v.minLength(3, "Nimessä tulee olla vähintään 3 merkkiä"),
						v.maxLength(50, "Nimessä voi olla enintään 50 merkkiä"),
				),
				email: v.pipe(
						v.string(),
						v.minLength(3, "Sähköpostin tulee olla vähintään 3 merkkiä"),
						v.email("Sähköposti on väärin formatoitu"),
						v.maxLength(50, "Sähköpostissa voi olla enintään 50 merkkiä"),
				),
				color: v.pipe(
						v.string(),
						v.nonEmpty("Valitse väri"),
						v.hexColor("???"),
				),
				admin: v.boolean()
		})
);

function UserForm(props) {
		const user = props.user;
		const cancel = () => props.cancel;

		const handleSubmit = async (data) => {
				await props.onSubmit(data);
		};

		const userForm = createForm({
				schema: userSchema,
		});

		createEffect(() => {
				if (user) {
						reset(userForm, {
								initialInput: user(),
						});
				}
		});

		return (
				<Form
						of={userForm}
						onSubmit={handleSubmit}
						class="form"
				>				
						<Field
								of={userForm}
								path={["name"]}
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
								of={userForm}
								path={["email"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Sähköposti:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="email"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						
						<Field
								of={userForm}
								path={["color"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Väri:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="color"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						
						<Field
								of={userForm}
								path={["admin"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Yllapito oikeudet:</label>
												<br/>
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="checkbox"
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

export default UserForm;