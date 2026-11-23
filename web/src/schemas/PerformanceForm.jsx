import { createForm, Field, Form, reset, setInput, getInput } from "@formisch/solid";
import * as v from "valibot";

const performanceSchema = v.pipe(
		v.object({
				points: v.pipe(
						v.number("Kirjoita numero"),
						v.integer("Kirjoita kokonaisluku"),
						v.toMinValue(1)
				),
		})
);

function PerformanceForm(props) {
		const cancel = () => props.cancel;

		const handleSubmit = async (data) => {
				await props.onSubmit(data);
		};

		const performanceForm = createForm({
				schema: performanceSchema,
				initialInput: {
						points: 1
				}
		});

		return (
				<Form
						of={performanceForm}
						onSubmit={handleSubmit}
				>				
            <div class="form-block">
                <strong>{props.title}</strong><br/>
                <span>{props.description}</span>
            </div>
            <div class="form-block">
								<Field
										of={performanceForm}
										path={["points"]}
								>
										{(field) => (
												<div>
														<input
																{...field.props}
																id={field.props.name}
																value={field.input}
																type="number"
																required
																onInput={(event)=>
																		setInput(performanceForm, {
																				path: ["points"],
																				input: parseInt(event.currentTarget.value)
																		})
																}
														/>
														<span>{props.unit}</span>
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
                    Vahvista
                </button>
            </div>
				</Form>
		);
}

export default PerformanceForm;