'use client'
import { Button } from "@heroui/button";
import {Form} from "@heroui/form";
import { Input } from "@heroui/input";
import coverImg from "../../assets/register.png"
import React from "react";
export default function Register() {
    const [action, setAction] = React.useState(null);
  return (
      <div className="container mx-auto">
        {/* wrapper */}
        <div className="flex flex-col lg:flex-row w-10/12 lg:w-8/12 bg-black-300 rounded-xl mx-auto shadow-lg overflow-hidden">
          {/* left */}
          <div
            className="w-full lg:w-1/2 flex flex-col items-center justify-center p-12 bg-no-repeat bg-center bg-cover"
            style={{
              backgroundImage: `url(${coverImg.src})`,
            }}
          >
         <h1 className="text-white text-5xl font-bold mb-3 text-shadow-lg">
  Welcome
</h1>
          </div>
          {/* right */}
          <div className="w-full lg:w-1/2 py-16 px-12">
            <h2 className="text-3xl mb-4">Register</h2>
            <p className="mb-4">
              Create your account. It’s free and only takes a minute
            </p>
            <Form
      className="w-full max-w-xs flex flex-col gap-4"
      onReset={() => {
        setAction("reset");
        console.log('Form reset');
      }}
      onSubmit={async (e) => {
        e.preventDefault();
        let data = Object.fromEntries(new FormData(e.currentTarget));
        console.log('Form data being sent:', data);
        setAction(`submitting...`);
        try {
          Console.log("Printing Data :" + data.stringify)
          const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
          });
          if (response.ok) {
            const result = await response.json();
            console.log('API response:', result);
            setAction(`submit success: ${JSON.stringify(result)}`);
          } else {
            const error = await response.json();
            console.log('API error:', error);
            setAction(`submit error: ${JSON.stringify(error)}`);
          }
        } catch (error) {
          console.error('Fetch error:', error);
          setAction(`error: ${error.message}`);
        }
      }}
    >
        <div className="flex row gap-3">
        <Input
        isRequired
        errorMessage="Please enter a valid username"
        label="First Name"
        labelPlacement="outside"
        name="firstname"
        placeholder="Enter your first name"
        type="text"
      />
      <Input
        isRequired
        errorMessage="Please enter a valid username"
        label="Last Name"
        labelPlacement="outside"
        name="lastname"
        placeholder="Enter your last name"
        type="text"
      /> </div>
     <Input
        isRequired
        errorMessage="Please enter a valid email"
        label="Email"
        labelPlacement="outside"
        name="email"
        placeholder="Enter your email"
        type="text"
      />
      <Input
        isRequired
        errorMessage="Please enter a valid password"
        label="Password"
        labelPlacement="outside"
        name="password"
        placeholder="Enter your password"
        type="password"
      />
      <Input
        isRequired
        errorMessage="Please enter a password that matches your first one"
        label="Confirm Password"
        labelPlacement="outside"
        name="confirm_password"
        placeholder="Confirm your password"
        type="password"
      />
      <div className="flex gap-2">
        <Button color="danger" size="lg" type="submit">
          Submit
        </Button>
        <Button type="reset" size="lg" variant="flat">
          Reset
        </Button>
      </div>
      {action && (
        <div className="text-small text-default-500">
          Action: <code>{action}</code>
        </div>
      )}
    </Form>
          </div>
        </div>
      </div>
  );
}