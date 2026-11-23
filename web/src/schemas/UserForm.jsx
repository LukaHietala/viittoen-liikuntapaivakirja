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
						v.hexColor("Väri ei ole oikeassa formaatissa"),
				),
				isAdmin: v.optional(v.boolean(), false),
		})
);

function UserForm(props) {
		const user = props.user;
		const cancel = () => props.cancel;
        const title = props.title
        const description = props.description

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
            <div class="form-block">
                <span>{title}</span><br/>
                <p class="small">{description}</p>
            </div>
            <div class="form-block">
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
										path={["isAdmin"]}
								>
										{(field) => (
												<div>
														<label for={field.props.name}>Yllapitäjä?: </label>
														<input
																{...field.props}
																id={field.props.name}
																checked={field.input}
																type="checkbox"
																required
																disabled={props.self}
														/>
														{field.errors && <small class="form-error">{field.errors[0]}</small>}
												</div>
										)}
								</Field>
            </div>
						
            <div class="form-actions">
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
            </div>
				</Form>
		);
}

export default UserForm;