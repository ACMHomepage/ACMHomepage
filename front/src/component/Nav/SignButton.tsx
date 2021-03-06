import { LogIn, LogOut } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';

import { useSelector } from '../../hooks';
import { selectAuthState, AuthStateEnum } from '../../store/authSlice';
import { URL as signInUrl } from '../../page/SignIn';
import { useSignOut } from '../../api/auth';

import LinkButton from '../LinkButton';

import styles from './styles/SignButton.module.scss';

/**
 * SignInOrRegisterButton. A button to go to sign in / register page.
 */
const SignInOrRegisterButton = (props: SignInOrRegisterProps) => {
  return (
    <LinkButton className={styles.button} to={signInUrl}>
      Sign in / Register
    </LinkButton>
  );
};

interface SignOutProps {
  className?: string;
}

/**
 * SignOutButton. A button to go to sign out.
 *
 * @param props - Holdes:
 * - className. `<string>`. To set its class.
 */
const SignOutButton = (props: SignOutProps) => {
  const { className } = props;

  const signOut = useSignOut();
  const authState = useSelector(selectAuthState);
  const navigate = useNavigate();

  useEffect(() => {
    if (authState === AuthStateEnum.Unlogged) {
      // If sign out, then go to the homepage.
      navigate('/');
    }
  }, [authState]);

  // TODO: Let it can sign out.
  return (
    <button className={styles.button} onClick={signOut}>
      Sign out
    </button>
  );
};

interface SignButtonProps {
  className?: string;
}

/**
 * SignButton. A button to handle the event of signing in, registering and so
 * on.
 *
 * @param props - Holdes:
 * - className. `<string>`. To set its class.
 */
const SignButton = (props: SignButtonProps) => {
  const { className } = props;
  const authState = useSelector(selectAuthState);

  return authState !== AuthStateEnum.LoggedWithInfo ? (
    <SignInOrRegisterButton className={className} />
  ) : (
    <SignOutButton className={className} />
  );
};

export default SignButton;
