package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.Utf8String;
import com.alaya.abi.solidity.datatypes.generated.Int256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class AbstractContractSon extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_FATHERAGE = "fatherAge";

    public static final String FUNC_FATHERNAME = "fatherName";

    public static final String FUNC_SONNAME = "sonName";

    protected AbstractContractSon(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AbstractContractSon(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> fatherAge() {
        final Function function = new Function(FUNC_FATHERAGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> fatherName() {
        final Function function = new Function(FUNC_FATHERNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> sonName() {
        final Function function = new Function(FUNC_SONNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<AbstractContractSon> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractSon.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AbstractContractSon> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractSon.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AbstractContractSon load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractSon(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AbstractContractSon load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractSon(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}