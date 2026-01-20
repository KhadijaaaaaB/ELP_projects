
import random
import string



def generate_random_UK_licence_plate(output_file="UK_SIV.txt"):
    letters1 = ''.join(random.choices(string.ascii_uppercase, k=2))
    digits = ''.join(random.choices(string.digits, k=2))
    letters2 = ''.join(random.choices(string.ascii_uppercase, k=3))
    licence_plate = f"{letters1} {digits} {letters2}"
    ## output licence plate into file
    with open(output_file, "a") as file:
        file.write(licence_plate + "\n")
    return licence_plate

def generate_random_french_licence_plate(output_file="French_SIV.txt"):

    letters = ''.join(random.choices(string.ascii_uppercase, k=2))
    digits = ''.join(random.choices(string.digits, k=3))
    region_code = ''.join(random.choices(string.ascii_uppercase + string.digits, k=2))
    licence_plate = f"{letters}-{digits}-{region_code}"
    ## output licence plate into file
    with open(output_file, "a") as file:
        file.write(licence_plate + "\n")
    return licence_plate

def generate_french_phone_number(output_file="french_phone_numbers.txt"):
    prefix = random.choice(['06', '07', '01', '02', '03', '04', '05', '09'])
    number = ''.join(random.choices(string.digits, k=8))
    phone_number = f"{prefix}{number}"
    ## output phone number into file
    with open(output_file, "a") as file:
        file.write(phone_number + "\n")
    return phone_number

def generate_emails_from_names(names_file="names.txt", output_file="emails.txt"):
    with open(names_file, "r") as file:
        names = [line.strip() for line in file.readlines()]
    emails = []
    for name in names:
        name_parts = name.lower().split()
        if len(name_parts) >= 2:
            first_name = name_parts[0]
            last_name = name_parts[-1]
            email = f"{first_name}.{last_name}@example.com"
            emails.append(email)
    with open(output_file, "w") as file:
        for email in emails:
            file.write(email + "\n")
    return emails

def shuffle_lists(list1, list2, list3=None, list4=None):
    # merge and shuffle two lists and write to file
    combined = list1 + list2
    if list3 is not None:
        combined += list3
    if list4 is not None:
        combined += list4
    random.shuffle(combined)
    with open("shuffled_output.txt", "w") as file:
        for item in combined:
            file.write(item + "\n")
    return combined

if __name__ == "__main__":
    for i in range(200):
        #generate_random_french_licence_plate()
        #generate_random_UK_licence_plate()
        #generate_french_phone_number()
        #generate_emails_from_names()
        pass

    with open("UK_SIV.txt", "r") as uk_file:
        list_uk = [line.strip() for line in uk_file.readlines()]
    with open("French_SIV.txt", "r") as fr_file:
        list_fr = [line.strip() for line in fr_file.readlines()]
    with open("french_phone_numbers.txt", "r") as fr_phone_file:
        list_fr_phones = [line.strip() for line in fr_phone_file.readlines()]
    with open("emails.txt", "r") as email_file:
        list_emails = [line.strip() for line in email_file.readlines()]
    shuffle_lists(list_uk, list_fr, list_fr_phones, list_emails)
    

    


    